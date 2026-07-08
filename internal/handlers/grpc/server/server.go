package server

import (
	"context"
	"database/sql"
	"fmt"

	auth_pb "github.com/LLM-Finetuning-platform/go-auth-api/internal/handlers/grpc/proto/v1"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/email"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/login"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/signup"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/utils"
	"github.com/redis/go-redis/v9"
	// import the generated protobuf code for the auth service
)

type MyServer struct{
	auth_pb.UnimplementedAuthAPIServiceServer
	Redis *redis.Client
	Resend *email.ResendClient
	Db *sql.DB
	From string
	Html string
	Subject string
	Secretconf *utils.TokenConfig
}

func (server *MyServer) Login (ctx context.Context, req *auth_pb.LoginRequest) (*auth_pb.LoginResponse, error){
	otp, err := utils.OTPGeneration()
	if err != nil{
		return nil, err
	}
	params := &email.EmailParams{From: server.From, To: []string{req.GetEmail().Email}, Html: server.Html, Subject: server.Subject}
	sreq := &email.Request{Email: req.GetEmail().Email, OTP: otp, Client: server.Resend, Params: params}
	err = login.Login(sreq, ctx, server.Redis)
	if err != nil {
		return  nil, err
	}
	status := &auth_pb.OTPStatus{OtpSent: true}
	return &auth_pb.LoginResponse{Status: status},nil

}

func (server *MyServer) LoginVerify (ctx context.Context, req *auth_pb.LoginVerifyRequest) (*auth_pb.LoginVerifyResponse, error){
	verify, err := login.LoginVerify(req.GetEmail().Email,  req.Request.GetOtp(), server.Redis, server.Db, ctx)
	if err != nil{
		return  nil, err
	}	

	username, err := utils.GetUsername(req.GetEmail().Email,ctx,server.Db)
	if err != nil{
		return nil , err
	}
	userdata := &auth_pb.UserData{Username: username, Email: req.GetEmail().Email}
	token, err := server.Secretconf.Encode(username,24,req.GetEmail().Email)
	if err != nil {
		return nil, err
	}
	otpres := &auth_pb.OTPResponse{Userdata: userdata, Token: &auth_pb.Token{Token: token}}
	return &auth_pb.LoginVerifyResponse{Status: verify, Response: otpres},nil
}

func (server *MyServer) Signup (ctx context.Context, req *auth_pb.SignupRequest) (*auth_pb.SignupResponse, error){
	otp, err := utils.OTPGeneration()
	if err != nil {
		return nil, err
	}
	params := &email.EmailParams{From: server.From, To: []string{req.GetEmail().Email}, Html: server.Html, Subject: server.Subject }	
	sreq := &email.Request{Email: req.Email.GetEmail(), OTP: otp, Client: server.Resend, Params: params}
	err = signup.SignUp(sreq, ctx ,server.Redis)
	if err != nil{
		return nil, err
	}
	status := &auth_pb.OTPStatus{OtpSent: true}
	return &auth_pb.SignupResponse{Status: status}, nil

}

func (server *MyServer) SignUPVerify (ctx context.Context, req *auth_pb.SignUPVerifyRequest) (*auth_pb.SignUPVerifyResponse, error) {
	verify, err := signup.SignUpVerify(ctx, req.GetEmail(),req.GetOtp(),req.GetUsername(),server.Redis,server.Db)
	if err!= nil {
		return  nil, err
	}
	userdata := &auth_pb.UserData{Username: req.GetUsername(),Email: req.GetEmail(),}
	token, err := server.Secretconf.Encode(req.GetUsername(), 24, req.GetEmail())
	if err!= nil{
		return nil, err
	}
	otpres := &auth_pb.OTPResponse{Userdata: userdata, Token: &auth_pb.Token{Token: token } }
	return &auth_pb.SignUPVerifyResponse{Status: verify, Response: otpres }, nil
		
}

func (server *MyServer) Auth (ctx context.Context, req *auth_pb.AuthRequest) (*auth_pb.AuthResponse, error) {
	usermap, err := server.Secretconf.Decode(req.GetToken().Token)		
	if err != nil {
		return nil, err
	}
	username, _ := usermap["username"].(string)
	email, _ := usermap["email"].(string)
	userdata := &auth_pb.UserData{Username: username,Email: email }
	return &auth_pb.AuthResponse{Userdata: userdata},nil
}

func  (server *MyServer) OTP (ctx context.Context, req *auth_pb.OTPRequest) (*auth_pb.OTPResponse, error) {
	verify, err := utils.Verify_otp(req.GetOtp(), ctx,req.Email.GetEmail(),server.Redis)
	if err !=nil {
		return nil, err
	}				
	if !verify {
		return nil, fmt.Errorf("OTP Not found")
	}
	username, err := utils.GetUsername(req.Email.GetEmail(),ctx,server.Db)
	if err !=nil{
		return nil, err
	} 	
	userdata := &auth_pb.UserData{Username: username, Email: req.Email.GetEmail()}
	token,err := server.Secretconf.Encode(username,24,req.Email.GetEmail())
	if err != nil  {
		return nil, err
	}
	return  &auth_pb.OTPResponse{Userdata: userdata,Token: &auth_pb.Token{Token: token}}, nil
}