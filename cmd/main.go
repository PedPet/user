package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/PedPet/proto/api/user"
	"github.com/PedPet/user/config"
	"github.com/PedPet/user/pkg/endpoint"
	userGrpc "github.com/PedPet/user/pkg/grpc"
	"github.com/PedPet/user/pkg/repository"
	"github.com/PedPet/user/pkg/service"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

const dbsource string = "root:Marthawinter1!@tcp(user-mysql)/pedpet"

func main() {
	grpcAddr := os.Getenv("PORT")

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "login",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var db *sql.DB
	{
		var err error

		db, err = sql.Open("mysql", dbsource)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	flag.Parse()
	ctx := context.Background()
	var cc service.CognitoClient
	{
		settings, err := config.LoadSettings()
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}

		conf := &aws.Config{
			Region:      aws.String("eu-west-1"),
			Credentials: credentials.NewStaticCredentials(settings.Aws.AccessKeyID, settings.Aws.SecretAccessKey, ""),
		}
		sess, err := session.NewSession(conf)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}

		identity := cognito.New(sess)
		cc, err = service.NewCognitoClient(identity, settings.Aws, logger)
	}

	var srv service.User
	{
		repository := repository.NewRepo(db, logger)
		srv = service.NewUserService(repository, cc, logger)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := endpoint.MakeEndpoints(srv)
	go func() {
		listener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			errs <- err
			return
		}

		handler := userGrpc.NewGRPCServer(ctx, endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterUserServer(gRPCServer, handler)
		errs <- gRPCServer.Serve(listener)
	}()

	level.Error(logger).Log("exit", <-errs)
}
