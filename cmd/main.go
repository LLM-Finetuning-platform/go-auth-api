package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/LLM-Finetuning-platform/go-auth-api/config"
	auth_pb "github.com/LLM-Finetuning-platform/go-auth-api/internal/handlers/grpc/proto/v2"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/handlers/grpc/server"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/db"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/email"
	redisservice "github.com/LLM-Finetuning-platform/go-auth-api/internal/services/redis"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/utils"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.GetNewConfig()
	if err != nil {

		fmt.Print(err)
		return
	}
	db_session, err := db.GetSession(cfg)
	if err != nil {
		fmt.Print(err)
		return
	}
	err = db_session.Ping()
	if err != nil {
		fmt.Print(err)
	}	
	stats := db_session.Stats()
	fmt.Printf("Successfully connected to the db, %v", stats.InUse)	
	redisservice := &redisservice.RedisService{Address: cfg.RedisAddress, Pwd: cfg.RedisPwd, DB: cfg.RedisDB, Protocol: cfg.RedisProtocol, Poolsize: cfg.RedisPoolsize, PoolTimeOutSecs: cfg.RedisPoolTimeOutSecs, MinIdleConnections: cfg.RedisMinIdleConnections}
	redisSession, err := redisservice.Connect()
	if err != nil {
		fmt.Print(err)
	}
	resend, err :=	email.NewClient(cfg)
	if err != nil {
		fmt.Print(err)
	}
	resendClient := &email.ResendClient{Createclient: email.NewClient, Client: resend}
	pemBytes, err := base64.StdEncoding.DecodeString(cfg.JwtPrivateKeyBase64)
	if err != nil{
		fmt.Print(err)
	}	
	secretConf, err := utils.NewTokenConfig(pemBytes)
	if err != nil {
		fmt.Print(err)
	}
	serverConf := &server.MyServer{Db: db_session,  Redis: redisSession, Resend: resendClient, Subject:    "Your Verification Code",
        Html:       "<p>Your OTP code is valid for 5 minutes.</p>",Secretconf: secretConf}
	listener, err:= net.Listen("tcp",cfg.Port)
	if err != nil {
		fmt.Print(err)
	}	
	grpcServer := grpc.NewServer()

	auth_pb.RegisterAuthAPIServiceServer(grpcServer,serverConf)
	stopChan := make(chan os.Signal,1)
	signal.Notify(stopChan,os.Interrupt, syscall.SIGTERM)
	fmt.Printf("gRPC server listening intently on port %s...\n",cfg.Port)
	go func(){
	if err = grpcServer.Serve(listener); err != nil{
		fmt.Print(err)
	}}()
	<-stopChan
	fmt.Println("\nShutdown signal received. Cleaning up resources...")

    grpcServer.GracefulStop()
    fmt.Println("gRPC server stopped gracefully.")
	if err = redisservice.Disconnect(redisSession); err != nil {
        fmt.Printf("Error closing Redis connection: %v\n", err)
    } else {
        fmt.Println("Redis pool disconnected cleanly.")
    }

    if err = db_session.Close(); err != nil {
        fmt.Printf("Error closing DB connection: %v\n", err)
    } else {
        fmt.Println("Database pool closed cleanly.")
    }

    fmt.Println("Shutdown complete. Exiting process.")
}