package service

import (
	"context"
	"os"
	"testing"

	"github.com/PedPet/user/config"
	"github.com/PedPet/user/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type needed struct {
	logger   log.Logger
	settings *config.Settings
	sess     *session.Session
	ctx      context.Context
	user     *model.User
}

func instantiateTest(t *testing.T) needed {
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

	settings, err := config.LoadSettings()
	if err != nil {
		t.Errorf("Failed to get settings: %v", err)
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

	ctx := context.Background()

	// user := &model.User{
	// 	Username:    "emmamartha",
	// 	Email:       "emmamarthacrossan@gmail.com",
	// 	Password:    "Swarleyfin1!",
	// 	PhoneNumber: "+447733814809",
	// }
	user := &model.User{
		Username:    "SC7639",
		Email:       "scrott@gmail.com",
		Password:    "Swarleyfin1!",
		PhoneNumber: "+447733814809",
	}

	return needed{
		logger:   logger,
		settings: settings,
		sess:     sess,
		ctx:      ctx,
		user:     user,
	}
}

func TestRegister(t *testing.T) {
	needed := instantiateTest(t)
	settings := needed.settings.Aws

	reg := &RegFlow{}
	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, reg, &settings, needed.logger)

	err = cc.Register(needed.ctx, needed.user)
	if err != nil {
		t.Errorf("Failed to register user: %v", err)
	}
}

func TestOTP(t *testing.T) {
	needed := instantiateTest(t)
	settings := needed.settings.Aws

	reg := &RegFlow{}
	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, reg, &settings, needed.logger)

	err = cc.OTP(needed.ctx, needed.user, "199967")
	if err != nil {
		t.Errorf("Failed to confirm OTP / user: %v", err)
	}
}

func TestResendConfirmation(t *testing.T) {
	needed := instantiateTest(t)
	settings := needed.settings.Aws

	reg := &RegFlow{}
	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, reg, &settings, needed.logger)
	if err != nil {
		t.Errorf("Failed to create new cognito client: %s", err)
	}
	err = cc.ResendConfirmation(needed.ctx, needed.user)
	if err != nil {
		t.Errorf("Failed to resend confirmation: %v", err)
	}
}

func TestCheckUsernameTaken(t *testing.T) {
	needed := instantiateTest(t)
	settings := needed.settings.Aws

	reg := &RegFlow{}
	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, reg, &settings, needed.logger)
	if err != nil {
		t.Errorf("Failed to create new cognito client: %s", err)
	}
	taken, err := cc.CheckUsernameTaken(needed.user.Username)
	if err != nil {
		t.Errorf("Failed to check if username is already taken: %v", err)
	}

	t.Logf("Username taken: %v", taken)
}

func TestLogin(t *testing.T) {
	needed := instantiateTest(t)
	settings := needed.settings.Aws

	reg := &RegFlow{}
	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, reg, &settings, needed.logger)
	if err != nil {
		t.Errorf("Failed to create new cognito client: %s", err)
	}
	err = cc.Login(needed.user.Username, needed.user.Password)
	if err != nil {
		t.Errorf("Failed to login: %v", err)
	}

}

func TestParseAndVerifyJWT(t *testing.T) {
	needed := instantiateTest(t)
	settings := needed.settings.Aws

	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, settings, needed.logger)
	if err != nil {
		t.Errorf("Failed to create new cognito client: %s", err)
	}

	_, err = cc.ParseAndVerifyJWT(needed.ctx, "eyJraWQiOiJtSlpDZ21EbDVzRUNvWUtTc1FYVk8wcFF5YkZJS2RjRks4aSswdXNcLytmZz0iLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJiM2EzYWMyZi1mN2YwLTQwODktOTczNy0wNmUzNTc5ZGYzYTIiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfMHBrZ051bzFGIiwicGhvbmVfbnVtYmVyX3ZlcmlmaWVkIjpmYWxzZSwiY3VzdG9tOnVpZCI6IjM3YjBiNWE5LTg5NTktNGQzOC1iYzkwLWVkNTU3OTEwYzkyZiIsImNvZ25pdG86dXNlcm5hbWUiOiJTQzc2MzkiLCJhdWQiOiI0djVvOW5jZHE0NDUxc3IwMGppbzdzbWwwaCIsImV2ZW50X2lkIjoiZmZmNjZhZWYtZmFmNi00NzYyLWFkMzgtYTc2ODQzMzc0YzU3IiwidG9rZW5fdXNlIjoiaWQiLCJhdXRoX3RpbWUiOjE1ODc3NzM2ODIsInBob25lX251bWJlciI6Iis0NDc3MzM4MTQ4MDkiLCJleHAiOjE1ODc3NzcyODIsImlhdCI6MTU4Nzc3MzY4MiwiZW1haWwiOiJzY3JvdHRAZ21haWwuY29tIn0.c_V0c7kbfViuL_HFeqMO4soAqcCL7KjnZaRaxmr3FaIw00lrSP_uZXzqZnhebj0G1WzgV9aC6YlHskC3t3dbdN24U4YxrVCBEAn2pIL9NomIMIJJyl0rEtvMOcznYhKaM3D6Op97Xaae3I15mh5J47InV543Re765F-U5xG3CvMTbZo5WLJwpDfqNAzMnhKWTgbQs456hhSrvm-TYN07x_PVSfhnJoFR7LieLTTW7_xpx1S1x1PfKb4NUtdVSMHzf8wInqLAeJsaFv7V7-gUZnS_i-64P-IT-dETBprqN1GtHMR9akcSiYeDwFrwIG-S_Yo-rYu5LD6OgmvA08maVA")
	if err != nil {
		t.Errorf("Failed to parse and verify jwt: %s", err)
	}
}

func TestGetUserDetails(t *testing.T) {
	needed := instantiateTest(t)
	settings := needed.settings.Aws

	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, settings, needed.logger)
	if err != nil {
		t.Errorf("Failed to create new cognito client: %s", err)
	}

	_, err = cc.GetUserDetails(needed.ctx, "eyJraWQiOiJNTzd5UGRpRXZoMUxkb2VDQmlrTmkxUXh1UG5ybEdPNEtHQURlaThheGl3PSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJiM2EzYWMyZi1mN2YwLTQwODktOTczNy0wNmUzNTc5ZGYzYTIiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xX2FlZTc0MzlhLWZiYjQtNDMzZC1hN2U1LTg1ODNmNWI0ZmVkZiIsImV2ZW50X2lkIjoiMjNkODhmMWItMTViZC00MzhiLWEyZWYtMmFlYzYxOTZiY2VmIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU4Nzg0NDczMiwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfMHBrZ051bzFGIiwiZXhwIjoxNTg3ODQ4MzMyLCJpYXQiOjE1ODc4NDQ3MzIsImp0aSI6IjIxZTJjNjg5LTg5ODktNDNlNC05ZWQwLTNiNTZlZWIzZmQwZSIsImNsaWVudF9pZCI6IjR2NW85bmNkcTQ0NTFzcjAwamlvN3NtbDBoIiwidXNlcm5hbWUiOiJTQzc2MzkifQ.CSLxe3_UnzZlSwSLh2KeNJDZxE8tCxozNCiGNbcaAXcLa6TTvMP873BGIy0_iZoprTCcDjJ_yRtEFmfElY5i3OTQZ7rxoLLjKIsjKzQNG0Tjdua3GLCmOPicjLkU1a4qHZZWG8VYzROpb8dfss2SDQJKX0HmWEcwdWhnEOUdjJyHtbl_cC3_oi4Phf8F5h4qEohpGSaU0pRYek786ymYhEcnZDZWksdTUgN4xxZUkmKrNJNsv5WP5PFnc2InabIVLpMGL6fO6QbS784q8BfzE9NwJ4R0GyTcFrv1dDRE11BDW9JCzhV3MobtSk7IHyVgoIsKkJi0BGTkbyAfmtiz8w", needed.user)
	if err != nil {
		t.Errorf("Failed to get user details: %s", err)
	}
}
