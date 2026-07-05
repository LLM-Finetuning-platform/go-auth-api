package server

import (
	"context"
	"database/sql"

	auth_pb "github.com/LLM-Finetuning-platform/go-auth-api/internal/handlers/grpc/proto/v1"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/email"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/login"
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
	secretconf *utils.TokenConfig
}

func (server *MyServer) Login (ctx context.Context, req *auth_pb.LoginRequest) (*auth_pb.LoginResponse, error){
	otp, err := utils.OTPGeneration()
	if err != nil{
		return nil, err
	}
	params := &email.EmailParams{From: server.from, To: []string{req.GetEmail().Email}, Html: server.html, Subject: server.subject}
	sreq := &email.Request{Email: req.GetEmail().Email, OTP: otp, Client: server.resend, Params: params}
	err = login.Login(sreq, ctx, server.redis)
	if err != nil {
		return  nil, err
	}
	status := &auth_pb.OTPStatus{OtpSent: true}
	return &auth_pb.LoginResponse{Status: status},nil

}

func (server *MyServer) LoginVerify (ctx context.Context, req *auth_pb.LoginVerifyRequest) (*auth_pb.LoginVerifyResponse, error){
	verify, err := login.LoginVerify(req.GetEmail().Email,  req.Request.GetOtp(), server.redis, server.db, ctx)
	if err != nil{
		return  nil, err
	}	
	userdata := &auth_pb.UserData{Username: req.}
}

func (server *MyServer) Signup (ctx context.Context, req *auth_pb.SignupRequest) (*auth_pb.SignupResponse, error){
	otp, err := utils.OTPGeneration()
	if err != nil {
		return nil, err
	}
	params := &email.EmailParams{From: server.from, To: []string{req.GetEmail().Email}, Html: server.html, Subject: server.subject }	
	sreq := &email.Request{Email: req.Email.GetEmail(), OTP: otp, Client: server.resend, Params: params}
	err = signup.SignUp(sreq, ctx ,server.redis)
	if err != nil{
		return nil, err
	}
	status := &auth_pb.OTPStatus{OtpSent: true}
	return &auth_pb.SignupResponse{Status: status}, nil

}

func (server *MyServer) SignUPVerify (ctx context.Context, req *auth_pb.SignUPVerifyRequest) (*auth_pb.SignUPVerifyResponse, error) {
	verify, err := signup.SignUpVerify(ctx, req.GetEmail(),req.GetOtp(),req.GetUsername(),server.redis,server.db)
	if err!= nil {
		return  nil, err
	}
	userdata := &auth_pb.UserData{Username: req.GetUsername(),Email: req.GetEmail(),}
	token, err := server.secretconf.Encode(req.GetUsername(), 24, req.GetEmail())
	if err!= nil{
		return nil, err
	}
	otpres := &auth_pb.OTPResponse{Userdata: userdata, Token: &auth_pb.Token{Token: token } }
	res := &auth_pb.SignUPVerifyResponse{Status: verify, Response: otpres }
	return res, nil 	
}