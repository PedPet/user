package grpc

import (
	"context"

	pb "github.com/PedPet/proto/api/user"
	"github.com/PedPet/user/pkg/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	createUser         grpctransport.Handler
	confirmUser        grpctransport.Handler
	resendConfirmation grpctransport.Handler
	usernameTaken      grpctransport.Handler
	login              grpctransport.Handler
	verifyJWT          grpctransport.Handler
	userDetails        grpctransport.Handler
}

// NewGRPCServer creates new user service
func NewGRPCServer(ctx context.Context, e endpoint.Endpoints) pb.UserServer {
	return &grpcServer{
		createUser: grpctransport.NewServer(
			e.CreateUser,
			endpoint.DecodeCreateUserRequest,
			endpoint.EncodeConfirmResponse,
		),
		confirmUser: grpctransport.NewServer(
			e.ConfirmUser,
			endpoint.DecodeConfirmUserRequest,
			endpoint.EncodeConfirmResponse,
		),
		resendConfirmation: grpctransport.NewServer(
			e.ResendConfirmation,
			endpoint.DecodeResendConfirmationRequest,
			endpoint.EncodeConfirmResponse,
		),
		usernameTaken: grpctransport.NewServer(
			e.UsernameTaken,
			endpoint.DecodeUsernameTakenRequest,
			endpoint.EncodeConfirmResponse,
		),
		login: grpctransport.NewServer(
			e.Login,
			endpoint.DecodeLoginRequest,
			endpoint.EncodeLoginResponse,
		),
		verifyJWT: grpctransport.NewServer(
			e.VerifyJWT,
			endpoint.DecodeVerifyJWTRequest,
			endpoint.EncodeConfirmResponse,
		),
		userDetails: grpctransport.NewServer(
			e.UserDetails,
			endpoint.DecodeUserDetailsRequest,
			endpoint.EncodeUserDetailsResponse,
		),
	}
}

func (s *grpcServer) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.ConfirmResponse, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.ConfirmResponse), nil
}

func (s *grpcServer) ConfirmUser(ctx context.Context, r *pb.ConfirmUserRequest) (*pb.ConfirmResponse, error) {
	_, resp, err := s.confirmUser.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.ConfirmResponse), nil
}

func (s *grpcServer) ResendConfirmation(ctx context.Context, r *pb.ResendConfirmationRequest) (*pb.ConfirmResponse, error) {
	_, resp, err := s.resendConfirmation.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.ConfirmResponse), nil
}

func (s grpcServer) UsernameTaken(ctx context.Context, r *pb.UsernameTakenRequest) (*pb.ConfirmResponse, error) {
	_, resp, err := s.usernameTaken.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.ConfirmResponse), nil
}

func (s *grpcServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, resp, err := s.confirmUser.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.LoginResponse), nil
}

func (s grpcServer) VerifyJWT(ctx context.Context, r *pb.VerifyJWTRequest) (*pb.ConfirmResponse, error) {
	_, resp, err := s.verifyJWT.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.ConfirmResponse), nil
}

func (s grpcServer) UserDetails(ctx context.Context, r *pb.UserDetailsRequest) (*pb.UserDetailsResponse, error) {
	_, resp, err := s.userDetails.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.UserDetailsResponse), nil
}
