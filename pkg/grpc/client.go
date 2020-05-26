package grpc

import (
	pb "github.com/PedPet/proto/api/user"
	"github.com/PedPet/user/pkg/endpoint"
	"github.com/PedPet/user/pkg/service"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func NewClient(conn *grpc.ClientConn) service.User {
	return endpoint.Endpoints{
		CreateUserEndpoint: grpctransport.NewClient(
			conn,
			"User",
			"CreateUser",
			endpoint.EncodeCreateUserRequest,
			endpoint.DecodeConfirmResponse,
			pb.ConfirmResponse{},
		).Endpoint(),
		ConfirmUserEndpoint: grpctransport.NewClient(
			conn,
			"User",
			"ConfirmUser",
			endpoint.EncodeConfirmUserRequest,
			endpoint.DecodeConfirmResponse,
			pb.ConfirmResponse{},
		).Endpoint(),
		ResendConfirmationEndpoint: grpctransport.NewClient(
			conn,
			"User",
			"ResendConfirmation",
			endpoint.EncodeResendConfirmationRequest,
			endpoint.DecodeConfirmResponse,
			pb.ConfirmResponse{},
		).Endpoint(),
		UsernameTakenEndpoint: grpctransport.NewClient(
			conn,
			"User",
			"UsernameTaken",
			endpoint.EncodeUsernameTakenRequest,
			endpoint.DecodeConfirmResponse,
			pb.ConfirmResponse{},
		).Endpoint(),
		LoginEndpoint: grpctransport.NewClient(
			conn,
			"User",
			"Login",
			endpoint.EncodeLoginRequest,
			endpoint.DecodeLoginResponse,
			pb.LoginResponse{},
		).Endpoint(),
		VerifyJWTEndpoint: grpctransport.NewClient(
			conn,
			"User",
			"VerifyJWT",
			endpoint.EncodeVerifyJWTRequest,
			endpoint.DecodeConfirmResponse,
			pb.ConfirmResponse{},
		).Endpoint(),
		UserDetailsEndpoint: grpctransport.NewClient(
			conn,
			"User",
			"UserDetails",
			endpoint.EncodeUserDetailsRequest,
			endpoint.DecodeUserDetailsResponse,
			pb.UserDetailsResponse{},
		).Endpoint(),
	}
}
