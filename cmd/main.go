package main

import (
	"fmt"

	"github.com/LLM-Finetuning-platform/go-auth-api/config"
	auth_pb "github.com/LLM-Finetuning-platform/go-auth-api/internal/handlers/grpc/proto/v1"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/handlers/grpc/server"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/db"
	"github.com/LLM-Finetuning-platform/go-auth-api/internal/services/email"
	redisservice "github.com/LLM-Finetuning-platform/go-auth-api/internal/services/redis"
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
	auth_server := &auth_pb.UnimplementedAuthAPIServiceServer{}
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
	serverConf := &server.MyServer{Db: db_session, Server: auth_server, Redis: redisSession, Resend: resendClient, Subject: "", From: "", Html: "",}
	err = redisservice.Disconnect(redisSession)
	if err != nil {
		fmt.Print(err)
	}
	err = db_session.Close()
	if err != nil {
		fmt.Print(err)
	}
}