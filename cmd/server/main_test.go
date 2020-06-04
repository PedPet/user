package main

import (
	"context"
	logg "log"
	"net"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/PedPet/proto/api/user"
	"github.com/PedPet/user/config"
	"github.com/PedPet/user/pkg/endpoint"
	grpcClient "github.com/PedPet/user/pkg/grpc"
	userGrpc "github.com/PedPet/user/pkg/grpc"
	"github.com/PedPet/user/pkg/repository"
	"github.com/PedPet/user/pkg/service"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/bxcodec/faker/v3"
	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type LoggerStub struct {
}

func (l LoggerStub) Log(keyvals ...interface{}) error {
	return nil
}

var lis *bufconn.Listener
var ctx context.Context
var mockGBL sqlmock.Sqlmock
var username string = faker.Username()
var password string = faker.Password() + "1!"
var jwt string

func init() {
	ctx = context.Background()
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	db, mock, err := sqlmock.New()
	if err != nil {
		logg.Fatalf("Failed to create stub database: %v", err)
	}
	// defer db.Close()
	mockGBL = mock

	// Instantiate logger
	var logger log.Logger
	{
		logger = LoggerStub{}
	}

	settings, err := config.LoadSettings()
	if err != nil {
		logg.Fatalf("Failed to load settings: %v", err)
	}

	// Instantiate conginto client
	var cc service.CognitoClient
	{
		conf := &aws.Config{
			Region:      aws.String("eu-west-1"),
			Credentials: credentials.NewStaticCredentials(settings.Aws.AccessKeyID, settings.Aws.SecretAccessKey, ""),
		}
		sess, err := session.NewSession(conf)
		if err != nil {
			logg.Fatalf("Failed to create new sessions: %v", err)
		}

		identity := cognito.New(sess)
		cc, err = service.NewCognitoClient(identity, settings.Aws, logger)
		if err != nil {
			logg.Fatalf("Failed to create cognito identity: %v", err)
		}
	}

	// Instantiate service
	var srv service.User
	{
		repository := repository.NewRepo(db, logger)
		srv = service.NewUserService(repository, cc, logger)
	}

	go func() {
		endpoints := endpoint.MakeEndpoints(srv)
		handler := userGrpc.NewGRPCServer(ctx, endpoints)
		pb.RegisterUserServer(s, handler)

		if err := s.Serve(lis); err != nil {
			logg.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(_ context.Context, _ string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "User", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	svc := grpcClient.NewClient(conn)

	mockGBL.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs(username).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = svc.CreateUser(ctx, username, "scrott@gmail.com", password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
}

func TestConfirmUser(t *testing.T) {
	t.Skip("Has to be set up from pervious test")
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "User", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	svc := grpcClient.NewClient(conn)
	err = svc.ConfirmUser(ctx, username, "")
	if err != nil {
		t.Fatalf("Failed to confirm user: %v", err)
	}
}

func TestResendConfirmation(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "User", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	svc := grpcClient.NewClient(conn)
	err = svc.ResendConfirmation(ctx, username)
	if err != nil {
		t.Fatalf("Failed to resend confirmation: %v", err)
	}
}

func TestUsernameTaken(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "User", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	svc := grpcClient.NewClient(conn)
	taken, err := svc.UsernameTaken(ctx, username)
	if err != nil {
		t.Fatalf("Failed to check if username is taken: %v", err)
	}

	if taken {
		t.Logf("Username is taken. Username: %s", username)
	} else {
		t.Logf("Username isn't taken. Username: %s", username)
	}
}

func TestLogin(t *testing.T) {
	// t.Skip("Needs to be set up with a confirmed user")

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "User", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %s", err)
	}
	defer conn.Close()

	svc := grpcClient.NewClient(conn)
	jwt, err = svc.Login(ctx, username, password)
	if err != nil {
		if strings.Contains(err.Error(), "UserNotConfirmedException") {
			t.Fatal("User is not confirmed")
		}
		t.Fatalf("Failed to login: %s", err)
	}

	t.Logf("JWT: %s", jwt)
}

func TestVerifyJWT(t *testing.T) {
	t.Skip("Needs to be set up with a JWT")

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "User", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %s", err)
	}
	defer conn.Close()

	svc := grpcClient.NewClient(conn)
	verified, err := svc.VerifyJWT(ctx, jwt)
	if err != nil {
		t.Fatalf("Failed to verify JWT: %s", err)
	}

	t.Logf("Verified: %v", verified)
}

func TestUserDetails(t *testing.T) {
	t.Skip("Needs to be set up with a JWT")

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "User", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %s", err)
	}
	defer conn.Close()

	if jwt == "" {
		t.Fatalf("JWT is empty: %s", jwt)
	}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	mockGBL.ExpectQuery("SELECT").WillReturnRows(rows)

	svc := grpcClient.NewClient(conn)
	user, err := svc.UserDetails(ctx, jwt)
	if err != nil {
		t.Fatalf("Failed to get user details: %s", err)
	}

	t.Logf("User: %v", user)
}
