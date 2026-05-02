package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eswarashish/go-auth-api/config"

	_ "github.com/lib/pq"
)


type Session struct{

}

func GetSession (conf *config.Config) (*sql.DB, error){
	dsn := conf.GetDatabaseUrl()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed create a postgres session , %w", err)	
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil{
		return  nil, fmt.Errorf("failed to ping database connection: %w", err)
	}
	return db, nil
}