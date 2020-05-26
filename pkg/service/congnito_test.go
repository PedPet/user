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

	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, settings, needed.logger)

	err = cc.Register(needed.ctx, needed.user)
	if err != nil {
		t.Errorf("Failed to register user: %v", err)
	}
}

func TestOTP(t *testing.T) {
	needed := instantiateTest(t)
	settings := needed.settings.Aws

	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, settings, needed.logger)

	err = cc.OTP(needed.ctx, needed.user, "199967")
	if err != nil {
		t.Errorf("Failed to confirm OTP / user: %v", err)
	}
}

func TestResendConfirmation(t *testing.T) {
	needed := instantiateTest(t)
	settings := needed.settings.Aws

	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, settings, needed.logger)
	if err != nil {
		t.Errorf("Failed to create new cognito client: %s", err)
	}
	err = cc.ResendConfirmation(needed.ctx, needed.user.Username)
	if err != nil {
		t.Errorf("Failed to resend confirmation: %v", err)
	}
}

func TestCheckUsernameTaken(t *testing.T) {
	needed := instantiateTest(t)
	settings := needed.settings.Aws

	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, settings, needed.logger)
	if err != nil {
		t.Errorf("Failed to create new cognito client: %s", err)
	}
	taken, err := cc.CheckUsernameTaken(needed.ctx, needed.user.Username)
	if err != nil {
		t.Errorf("Failed to check if username is already taken: %v", err)
	}

	t.Logf("Username taken: %v", taken)
}

func TestLogin(t *testing.T) {
	needed := instantiateTest(t)
	settings := needed.settings.Aws

	identity := cognito.New(needed.sess)
	cc, err := NewCognitoClient(identity, settings, needed.logger)
	if err != nil {
		t.Errorf("Failed to create new cognito client: %s", err)
	}
	_, err = cc.Login(needed.ctx, needed.user.Username, needed.user.Password)
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

	_, err = cc.GetUserDetails(needed.ctx, "eyJraWQiOiJTSmVTM2VIT1pvSFFWaXErcUlcL1lQNzVpNjhnQldQQXJmcUEwRUI4ZnZZMD0iLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIzNTNhNzdlYi00NzI3LTQ5YTgtYjU1Yi0yODdkMGNjZTBjMjQiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xXzVhMmY0ZTdmLTFiN2QtNGVhNy05MzQ4LTQ0ODIwYTllZDhiMiIsImV2ZW50X2lkIjoiNjg4MTNmMjItNTgxMi00YjA0LTlhMDEtMzc2YjQyODg3Mjk2IiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5MDUwMzk4OSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfOHNyOG9EdVJLIiwiZXhwIjoxNTkwNTA3NTg5LCJpYXQiOjE1OTA1MDM5OTAsImp0aSI6IjQ4ZmMwYzdjLWNjZjAtNGZkNy05OGY2LTkyNDFiZTgxY2QwYiIsImNsaWVudF9pZCI6IjRrdXBqZGg4ZnEzdWJsMTc3YjZjMWZidWF0IiwidXNlcm5hbWUiOiJzYzc2MzkifQ.GTClCaUBf3Yi2fxZiT0Lg-HU0YoMmIlNIFI02TiPHBSCipGELw3ERaTZYPSQNd90FomcvRJlUVvbIVjGh_Z8DgAoBjXnafA0wJJh_MgeTul03vh6NAbCjc8H9AoiHj-ALuKv8Bs63WxBpXV0X6sXh72V0GGyzYbfsNV_DK9dAd7XbUyo4tqUlux3QK47IkhK0CjqjxmXAk5RDbCiYHl1yEEm-tTvq-1hJNXaluFZieUhdb8Tuo3OuAmUY6sEQEZZV3q_bxkjAWFL4vecFYWPBHALim94Clle4VIyIW5J1PcD46Ut0wIzJm9CEvQko3gmYkBI8Jtdi0P4dId9UyQkGQ")
	if err != nil {
		t.Errorf("Failed to get user details: %s", err)
	}
}
