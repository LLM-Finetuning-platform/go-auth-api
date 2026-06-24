package server

import (
	"context"

	auth_pb "github.com/LLM-Finetuning-platform/go-auth-api/internal/handlers/grpc/proto/v1"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/signup"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/utils"
	// import the generated protobuf code for the auth service
)

type MyServer struct{
	server *auth_pb.UnimplementedAuthAPIServiceServer
	

}

func (server *MyServer) Login (ctx context.Context, req *auth_pb.LoginRequest) (*auth_pb.LoginResponse, error){
	return &auth_pb.LoginResponse{}, nil
}

func (server *MyServer) Signup (ctx context.Context, req *auth_pb.SignupRequest) (*auth_pb.SignupResponse, error){
	otp, err := utils.OTPGeneration()
	if err != nil {
		return nil, err
	}
	sreq := &signup.SignUpRequest{Email: req.Email.GetEmail(), OTP: otp, Client: nil, Params: nil}
	err = signup.SignUp(sreq, nil)
	if err != nil{
		return nil, err
	}
	status := &auth_pb.OTPStatus{OtpSent: true}
	return &auth_pb.SignupResponse{Status: status}, nil

}
