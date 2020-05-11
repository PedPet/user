package service

import (
	"context"

	"github.com/PedPet/user/model"
	"github.com/PedPet/user/pkg/repository"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// User describes the service.
type User interface {
	CreateUser(ctx context.Context, username, email, password string) (*model.User, error)
	UserDetails(ctx context.Context, token string) (*model.User, error)
	ConfirmUser(ctx context.Context, username, otp string) (*model.User, error)
	ResendConfirmation(ctx context.Context, username string) error
	CheckUsernameTaken(ctx context.Context, username string) (bool, error)
	Login(ctx context.Context, username, password string) (string, error)
	VerifyJWT(ctx context.Context, token string) (bool, error)
}

type service struct {
	repository repository.User
	cognito    CognitoClient
	logger     log.Logger
}

// NewLoginService creates a login service with required dependencies
func NewUserService(rep repository.User, cognito CognitoClient, logger log.Logger) User {
	return &service{
		repository: rep,
		cognito:    cognito,
		logger:     logger,
	}
}

// CreateUser registers a user with cognito and then stores the username with an id for further user data storage
func (s service) CreateUser(ctx context.Context, username email, password string) (*model.User, error) {
	logger := log.With(s.logger, "method", "CreateUser")

	user := &model.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	err := s.cognito.Register(ctx, user)
	if err != nil {
		level.Error(logger).Log("err", err)
	}

	err = s.repository.CreateUser(ctx, user)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	logger.Log("Create user", user)
	return user, nil
}

// GetUser gets the user's details from cognito and the database and
func (s service) UserDetails(ctx context.Context, token string) (*model.User, error) {
	logger := log.With(s.logger, "method", "GetUser")

	user, err := s.cognito.GetUserDetails(ctx, token)

	err = s.repository.GetUser(ctx, user)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	logger.Log("User details", user)

	return user, nil
}

func (s service) ConfirmUser(ctx context.Context, username, otp string) (*model.User, error) {
	logger := log.With(s.logger, "method", "ConfirmUser")

	user := &model.User{
		Username: username,
	}
	err := s.cognito.OTP(ctx, user, otp)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	user.Confirmed = true

	logger.Log("Confirm user", user)
	return user, nil
}

func (s service) ResendConfirmation(ctx context.Context, username string) error {
	logger := log.With(s.logger, "method", "ResendConfirmation")

	err := s.cognito.ResendConfirmation(ctx, username)
	if err != nil {
		level.Error(logger).Log("err", err)
		return err
	}

	logger.Log("Resend confirmation")
	return nil
}

func (s service) CheckUsernameTaken(ctx context.Context, username string) (bool, error) {
	logger := log.With(s.logger, "method", "CheckUsernameTaken")

	taken, err := s.cognito.CheckUsernameTaken(ctx, username)
	if err != nil {
		level.Error(logger).Log("err", err)
		return false, nil
	}

	logger.Log("Username taken", taken)
	return taken, nil
}

func (s service) Login(ctx context.Context, username, password string) (string, error) {
	logger := log.With(s.logger, "method", "Login")

	auth, err := s.cognito.Login(ctx, username, password)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("Login", auth)
	return aws.StringValue(auth.IdToken), nil
}

func (s service) VerifyJWT(ctx context.Context, jwt string) (bool, error) {
	logger := log.With(s.logger, "method", "VerifyJWT")

	token, err := s.cognito.ParseAndVerifyJWT(ctx, jwt)
	if err != nil {
		level.Error(logger).Log("err", err)
		return false, err
	}

	return token.Valid, nil
}
