package main

import (
	"fmt"

	"github.com/eswarashish/go-auth-api/config"
	"github.com/eswarashish/go-auth-api/internal/services/db"
)

func main()  {
	cfg, err := config.GetNewConfig()
	if err !=nil{
		
		fmt.Print(err)
		return
	}
	session, err := db.GetSession(cfg)
	if err != nil{
		fmt.Print(err)
		return
	}
	err = session.Ping()
	if err != nil {
		fmt.Print(err)
	}
	stats := session.Stats()
	fmt.Printf("Successfully connected to the db, %v", stats.InUse)
}