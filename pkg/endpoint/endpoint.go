package endpoint

import (
	"context"

	service "github.com/PedPet/user/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

// Endpoints is a struct that contains all the endpoint available in this microservice
type Endpoints struct {
	CreateUser         endpoint.Endpoint
	ConfirmUser        endpoint.Endpoint
	ResendConfirmation endpoint.Endpoint
	UsernameTaken      endpoint.Endpoint
	Login              endpoint.Endpoint
	VerifyJWT          endpoint.Endpoint
	UserDetails        endpoint.Endpoint
}

// MakeEndpoints give the required dependencies to the Endpoints
func MakeEndpoints(s service.User) Endpoints {
	return Endpoints{
		CreateUser:         makeCreateUserEndpoint(s),
		ConfirmUser:        makeConfirmUser(s),
		ResendConfirmation: makeResendConfirmation(s),
		UsernameTaken:      makeUsernameTaken(s),
		Login:              makeLogin(s),
		VerifyJWT:          makeVerifyJWT(s),
		UserDetails:        makeUserDetails(s),
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

func makeConfirmUser(s service.User) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ConfirmUserRequest)
		_, err := s.ConfirmUser(ctx, req.Username, req.Code)
		if err != nil {
			return nil, err
		}

		return ConfirmResponse{Ok: true}, nil
	}
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

func makeUsernameTaken(s service.User) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UsernameTakenRequest)
		t, err := s.CheckUsernameTaken(ctx, req.Username)
		if err != nil {
			return nil, err
		}

		return ConfirmResponse{Ok: t}, nil
	}
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

func makeVerifyJWT(s service.User) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyJWTRequest)
		_, err := s.VerifyJWT(ctx, req.Jwt)
		if err != nil {
			return nil, err
		}

		return VerifyJWTRequest{Jwt: req.Jwt}, nil
	}
}

func makeUserDetails(s service.User) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UserDetailsRequest)
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
