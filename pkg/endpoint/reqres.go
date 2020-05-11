package endpoint

import (
	"context"

	pb "github.com/PedPet/proto/api/user"
)

type (
	// ConfirmResponse test
	ConfirmResponse struct {
		Ok bool `json:"ok"`
	}

	// CreateUserRequest is a struct to convert a create user request to and from json
	CreateUserRequest struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phoneNumber"`
	}

	// ConfirmUserRequest is a struct to convert a get user response to and from json
	ConfirmUserRequest struct {
		Username string `json:"username"`
		Code     string `json:"code"`
	}

	// ResendConfirmationRequest test
	ResendConfirmationRequest struct {
		Username string `json:"username"`
	}

	// UsernameTakenRequest test
	UsernameTakenRequest struct {
		Username string `json:"username"`
	}

	// LoginRequest test
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// LoginResponse test
	LoginResponse struct {
		Jwt string `json:"jwt"`
	}

	// VerifyJWTRequest test
	VerifyJWTRequest struct {
		Jwt string `json:"jwt"`
	}

	// UserDetailsRequest test
	UserDetailsRequest struct {
		Jwt string `json:"jwt"`
	}

	// UserDetailsResponse test
	UserDetailsResponse struct {
		ID          int    `json:"id"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
		Confirmed   bool   `json:"confirmed"`
	}
)

// EncodeConfirmResponse encode internal response into grpc response type
func EncodeConfirmResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(ConfirmResponse)
	return &pb.ConfirmResponse{
		Ok: resp.Ok,
	}, nil
}

// DecodeCreateUserRequest decode grpc request into expected internal request type
func DecodeCreateUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CreateUserRequest)
	return CreateUserRequest{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
	}, nil
}

// DecodeConfirmUserRequest decode grpc request into expected internal request type
func DecodeConfirmUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.ConfirmUserRequest)
	return ConfirmUserRequest{
		Username: req.Username,
		Code:     req.Code,
	}, nil
}

// DecodeResendConfirmationRequest decode grpc request into expected internal request type
func DecodeResendConfirmationRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.ResendConfirmationRequest)
	return ResendConfirmationRequest{
		Username: req.Username,
	}, nil
}

// DecodeUsernameTakenRequest decode grpc request into expected internal request type
func DecodeUsernameTakenRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UsernameTakenRequest)
	return UsernameTakenRequest{
		Username: req.Username,
	}, nil
}

// DecodeLoginRequest decode grpc request into expected internal request type
func DecodeLoginRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.LoginRequest)
	return LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}, nil
}

// EncodeLoginResponse encode the internal response into the expected grpc response type
func EncodeLoginResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(LoginResponse)
	return &pb.LoginResponse{
		Jwt: resp.Jwt,
	}, nil
}

// DecodeVerifyJWTRequest decode the grpc request into the expected internal request type
func DecodeVerifyJWTRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.VerifyJWTRequest)
	return &VerifyJWTRequest{
		Jwt: req.Jwt,
	}, nil
}

// DecodeUserDetailsRequest decode the grpc request into the expected internal request type
func DecodeUserDetailsRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UserDetailsRequest)
	return &UserDetailsRequest{
		Jwt: req.Jwt,
	}, nil
}

// EncodeUserDetailsResponse encode the internal response into the expected grpc response type
func EncodeUserDetailsResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(UserDetailsResponse)
	return &pb.UserDetailsResponse{
		Id:          int32(resp.ID),
		Username:    resp.Username,
		Email:       resp.Email,
		PhoneNumber: resp.PhoneNumber,
		Confirmed:   resp.Confirmed,
	}, nil
}
