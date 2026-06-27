package server

import (
	"context"
	"database/sql"

	auth_pb "github.com/LLM-Finetuning-platform/go-auth-api/internal/handlers/grpc/proto/v1"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/email"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/signup"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/utils"
	"github.com/redis/go-redis/v9"
	// import the generated protobuf code for the auth service
)

type MyServer struct{
	server *auth_pb.UnimplementedAuthAPIServiceServer
	redis *redis.Client
	resend *email.ResendClient
	db *sql.DB
	from string
	html string
	subject string
}

func (server *MyServer) Login (ctx context.Context, req *auth_pb.LoginRequest) (*auth_pb.LoginResponse, error){
	return &auth_pb.LoginResponse{}, nil
}

func (server *MyServer) Signup (ctx context.Context, req *auth_pb.SignupRequest) (*auth_pb.SignupResponse, error){
	otp, err := utils.OTPGeneration()
	if err != nil {
		return nil, err
	}
	params := &email.EmailParams{From: server.from, To: []string{req.GetEmail().Email}, Html: server.html, Subject: server.subject }	
	sreq := &signup.SignUpRequest{Email: req.Email.GetEmail(), OTP: otp, Client: server.resend, Params: params}
	err = signup.SignUp(sreq,  server.redis)
	if err != nil{
		return nil, err
	}
	status := &auth_pb.OTPStatus{OtpSent: true}
	return &auth_pb.SignupResponse{Status: status}, nil

}

