package endpoint

import (
	"context"

	"github.com/PedPet/user/model"
	service "github.com/PedPet/user/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)

// Endpoints is a struct that contains all the endpoint available in this microservice
type Endpoints struct {
	CreateUserEndpoint         endpoint.Endpoint
	ConfirmUserEndpoint        endpoint.Endpoint
	ResendConfirmationEndpoint endpoint.Endpoint
	UsernameTakenEndpoint      endpoint.Endpoint
	LoginEndpoint              endpoint.Endpoint
	VerifyJWTEndpoint          endpoint.Endpoint
	UserDetailsEndpoint        endpoint.Endpoint
}

// MakeEndpoints give the required dependencies to the Endpoints
func MakeEndpoints(s service.User) Endpoints {
	return Endpoints{
		CreateUserEndpoint:         makeCreateUserEndpoint(s),
		ConfirmUserEndpoint:        makeConfirmUser(s),
		ResendConfirmationEndpoint: makeResendConfirmation(s),
		UsernameTakenEndpoint:      makeUsernameTaken(s),
		LoginEndpoint:              makeLogin(s),
		VerifyJWTEndpoint:          makeVerifyJWT(s),
		UserDetailsEndpoint:        makeUserDetails(s),
	}
}

func makeCreateUserEndpoint(s service.User) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		_, err := s.CreateUser(ctx, req.Username, req.Email, req.Password)
		if err != nil {
			return nil, err
		}

		return ConfirmResponse{Ok: true}, nil
	}
}

// CreateUser calls the create user endpoint
func (e Endpoints) CreateUser(ctx context.Context, username, email, password string) (*model.User, error) {
	req := CreateUserRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	resp, err := e.CreateUserEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	createUserResp := resp.(ConfirmResponse)
	if createUserResp.Ok != true {
		return nil, errors.New("Failed to create user")
	}

	return &model.User{
		Username: username,
		Email:    email,
		Password: password,
	}, nil
}

func makeConfirmUser(s service.User) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ConfirmUserRequest)
		err := s.ConfirmUser(ctx, req.Username, req.Code)
		if err != nil {
			return nil, err
		}

		return ConfirmResponse{Ok: true}, nil
	}
}

// ConfirmUser calls the confirm user endpoint
func (e Endpoints) ConfirmUser(ctx context.Context, username, otp string) error {
	req := ConfirmUserRequest{
		Username: username,
		Code:     otp,
	}

	resp, err := e.ConfirmUserEndpoint(ctx, req)
	if err != nil {
		return err
	}

	confirmUserResp := resp.(ConfirmResponse)
	if confirmUserResp.Ok != true {
		return errors.New("Failed to confirm user")
	}
	return nil
}

func makeResendConfirmation(s service.User) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ResendConfirmationRequest)
		err := s.ResendConfirmation(ctx, req.Username)
		if err != nil {
			return nil, err
		}

		return ConfirmResponse{Ok: true}, nil
	}
}

// ResendConfirmation calls the resend confirmation endpoint
func (e Endpoints) ResendConfirmation(ctx context.Context, username string) error {
	req := ResendConfirmationRequest{
		Username: username,
	}

	resp, err := e.ResendConfirmationEndpoint(ctx, req)
	if err != nil {
		return err
	}

	resendConfirmationResp := resp.(ConfirmResponse)
	if resendConfirmationResp.Ok != true {
		return errors.New("Failed to resend confirmation")
	}
	return nil

}

func makeUsernameTaken(s service.User) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UsernameTakenRequest)
		t, err := s.UsernameTaken(ctx, req.Username)
		if err != nil {
			return nil, err
		}

		return ConfirmResponse{Ok: t}, nil
	}
}

// UsernameTaken calls the check username taken endpoint
func (e Endpoints) UsernameTaken(ctx context.Context, username string) (bool, error) {
	req := UsernameTakenRequest{
		Username: username,
	}

	resp, err := e.UsernameTakenEndpoint(ctx, req)
	if err != nil {
		return false, err
	}

	checkUsernameResp := resp.(ConfirmResponse)
	return checkUsernameResp.Ok, nil
}

func makeLogin(s service.User) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		token, err := s.Login(ctx, req.Username, req.Password)
		if err != nil {
			return nil, err
		}

		return LoginResponse{Jwt: token}, nil
	}
}

// Login calls the login endpoint
func (e Endpoints) Login(ctx context.Context, username, password string) (string, error) {
	req := LoginRequest{
		Username: username,
		Password: password,
	}

	resp, err := e.LoginEndpoint(ctx, req)
	if err != nil {
		return "", err
	}

	loginResp := resp.(LoginResponse)

	return loginResp.Jwt, nil
}

func makeVerifyJWT(s service.User) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyJWTRequest)
		verified, err := s.VerifyJWT(ctx, req.Jwt)
		if err != nil {
			return nil, err
		}

		if verified != true {
			return ConfirmResponse{Ok: false}, nil
		}

		return ConfirmResponse{Ok: true}, nil
	}
}

// VerifyJWT calls the verify jwt endpoint
func (e Endpoints) VerifyJWT(ctx context.Context, token string) (bool, error) {
	req := VerifyJWTRequest{
		Jwt: token,
	}

	resp, err := e.VerifyJWTEndpoint(ctx, req)
	if err != nil {
		return false, err
	}

	verifyJWTResp := resp.(ConfirmResponse)
	return verifyJWTResp.Ok, nil
}

func makeUserDetails(s service.User) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*UserDetailsRequest)
		user, err := s.UserDetails(ctx, req.Jwt)
		if err != nil {
			return nil, err
		}

		return UserDetailsResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Confirmed:   user.Confirmed,
		}, nil
	}
}

// UserDetails calls the user details enpoint
func (e Endpoints) UserDetails(ctx context.Context, token string) (*model.User, error) {
	req := UserDetailsRequest{
		Jwt: token,
	}

	resp, err := e.UserDetailsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	userDetailsResp := resp.(UserDetailsResponse)
	return &model.User{
		ID:          userDetailsResp.ID,
		Username:    userDetailsResp.Username,
		Email:       userDetailsResp.Email,
		PhoneNumber: userDetailsResp.PhoneNumber,
		Confirmed:   userDetailsResp.Confirmed,
	}, nil
}
