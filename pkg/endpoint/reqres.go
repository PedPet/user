package endpoint

import (
	"context"
	"regexp"

	pb "github.com/PedPet/proto/api/user"
	"github.com/PedPet/user/config"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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

// DecodeConfirmResponse decode grpc response into internal response
func DecodeConfirmResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.ConfirmResponse)
	return ConfirmResponse{
		Ok: resp.Ok,
	}, nil
}

// EncodeCreateUserRequest encode internal request into grpc request type
func EncodeCreateUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(CreateUserRequest)
	return &pb.CreateUserRequest{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
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

// EncodeConfirmUserRequest encodes internal request into a grpc request type
func EncodeConfirmUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(ConfirmUserRequest)
	return &pb.ConfirmUserRequest{
		Username: req.Username,
		Code:     req.Code,
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

// EncodeResendConfirmationRequest encodes internal request into grpc request type
func EncodeResendConfirmationRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(ResendConfirmationRequest)
	return &pb.ResendConfirmationRequest{
		Username: req.Username,
	}, nil
}

// DecodeResendConfirmationRequest decode grpc request into expected internal request type
func DecodeResendConfirmationRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.ResendConfirmationRequest)
	return ResendConfirmationRequest{
		Username: req.Username,
	}, nil
}

// EncodeUsernameTakenRequest encode request into expected grpc request
func EncodeUsernameTakenRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(UsernameTakenRequest)
	return &pb.UsernameTakenRequest{
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

// EncodeLoginRequest encode request into expected grpc request
func EncodeLoginRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(LoginRequest)
	return &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
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

// DecodeLoginResponse decode grpc response into expected internal response type
func DecodeLoginResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.LoginResponse)
	return LoginResponse{
		Jwt: resp.Jwt,
	}, nil
}

// EncodeVerifyJWTRequest encodes internal request into the expected grpc request type
func EncodeVerifyJWTRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(VerifyJWTRequest)
	return &pb.VerifyJWTRequest{
		Jwt: req.Jwt,
	}, nil
}

// DecodeVerifyJWTRequest decode the grpc request into the expected internal request type
func DecodeVerifyJWTRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.VerifyJWTRequest)
	return VerifyJWTRequest{
		Jwt: req.Jwt,
	}, nil
}

// EncodeUserDetailsRequest encodes the internal request into the grpc request type
func EncodeUserDetailsRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(UserDetailsRequest)
	return &pb.UserDetailsRequest{
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

// DecodeUserDetailsResponse seconds the grpc response into the expected internal response type
func DecodeUserDetailsResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.UserDetailsResponse)
	return UserDetailsResponse{
		ID:          int(resp.Id),
		Username:    resp.Username,
		Email:       resp.Email,
		PhoneNumber: resp.PhoneNumber,
		Confirmed:   resp.Confirmed,
	}, nil
}

// Username cannot be empty and must be between 2 and 100 characters
func validUsername(username *string) *validation.FieldRules {
	return validation.Field(username, validation.Required, validation.Length(2, 100))
}

// Password cannot be empty and must be between 8 and 100 characters
func validPassword(password *string, rules config.Password) *validation.FieldRules {
	rr := []validation.Rule{
		validation.Required,
		validation.Length(rules.Length.Min, rules.Length.Max),
	}
	for _, regex := range rules.Regex {
		rr = append(rr, validation.Match(regexp.MustCompile(regex)))
	}

	return validation.Field(
		password,
		rr...,
	)
}

// Validate the request payload
func (r CreateUserRequest) Validate(pwRules config.Password) error {
	return validation.ValidateStruct(&r,
		validUsername(&r.Username),
		// Email connot be empty and must be a valid email address
		validation.Field(&r.Email, validation.Required, is.Email),
		validPassword(&r.Password, pwRules),
	)
}

// Validate the request payload
func (r ConfirmUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validUsername(&r.Username),
		// Code cannot be empty and must be 6 digits
		validation.Field(&r.Code, validation.Required, is.Digit, validation.Length(6, 6)),
	)
}

// Validate the request payload
func (r ResendConfirmationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validUsername(&r.Username),
	)
}

// Validate the request payload
func (r UsernameTakenRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validUsername(&r.Username),
	)
}

// Validate the request payload
func (r LoginRequest) Validate(rules config.Password) error {
	return validation.ValidateStruct(&r,
		validUsername(&r.Username),
		validPassword(&r.Password, rules),
	)
}

// Validate the request payload
func (r VerifyJWTRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Jwt, validation.Required),
	)
}

// Validate the request payload
func (r UserDetailsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Jwt, validation.Required),
	)
}
