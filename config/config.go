package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct{
	DatabaseName string `env:"POSTGRES_DB"`
	DatabasePassword string `env:"POSTGRES_PASSWORD"`
	DatabaseUser string `env:"POSTGRES_USER"`
	DatabasePort string `env:"POSTGRES_PORT" envDefault:"5432"`	
	DatbaseHost string `env:"POSTGRES_HOST"`
}

//function attached as property to string by using reciever
func (conf *Config) GetDatabaseUrl() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",conf.DatabaseUser,conf.DatabasePassword,conf.DatbaseHost,conf.DatabasePort,conf.DatabaseName)
}

func GetNewConfig() (*Config, error) {
	//since this is only dev environment we are loading in dev
	err := godotenv.Load(".env.dev")
	if err !=nil{
		return nil, fmt.Errorf("Failed to load the dev environment %w", err)
	}
	cfg, err := env.ParseAs[Config]();
	if err != nil {
		return nil, fmt.Errorf("failed to load config %w",err)
	} 
	return &cfg, nil
}

