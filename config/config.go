package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/eswarashish/go-auth-api/internal/services/redisservice"
	"github.com/eswarashish/go-auth-api/internal/utils/logger"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct{

	DatabaseName string `env:"POSTGRES_DB"`
	DatabasePassword string `env:"POSTGRES_PASSWORD"`
	DatabaseUser string `env:"POSTGRES_USER"`
	DatabasePort string `env:"POSTGRES_PORT" envDefault:"5432"`	
	DatabaseHost string `env:"POSTGRES_HOST"`
}
type Config struct{
	ResendAPIKey string `env:"RESEND_API_KEY"`
	DBConfig DatabaseConfig
	CacheConfig redisservice.RedisService
	CacheTTL time.Duration
}

//function attached as property to string by using reciever
func (conf *Config) GetDatabaseUrl() string {
	logger := logger.AuthSlogger.GetLogger()
	logger.Debug("Fetching the Database Url as string")
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",conf.DBConfig.DatabaseUser,
	conf.DBConfig.DatabasePassword,
	conf.DBConfig.DatabaseHost,
	conf.DBConfig.DatabasePort,
	conf.DBConfig.DatabaseName)
}

func GetNewConfig() (*Config, error) {
	//since this is only dev environment we are loading in dev
	err := loadEnv()
	if err !=nil{
		return nil, fmt.Errorf("Failed to load the dev environment %w", err)
	}
	cfg, err := env.ParseAs[Config]();
	if err != nil {
		return nil, fmt.Errorf("failed to load config %w",err)
	} 
	return &cfg, nil
}

func loadEnv() (error) {
	appEnv := os.Getenv("env")
	if appEnv == ""{
		appEnv = "dev"
	}
	if (appEnv == "prod") || (appEnv == "staging"){
		return nil
	}
	envFile := fmt.Sprintf(".env.%s",appEnv)
	if err := godotenv.Load(envFile); err != nil{
		return fmt.Errorf("Failed to load %s: %w", envFile, err)
	}
	return nil
} 