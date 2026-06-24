package server

import (
	auth_pb "github.com/LLM-Finetuning-platform/go-auth-api/internal/handlers/grpc/proto/v1"
	
	// import the generated protobuf code for the auth service
)

type MyServer struct{
	server *auth_pb.UnimplementedAuthAPIServiceServer
	

}