package utils

import (
	"context"
	"database/sql"

	"github.com/LLM-Finetuning-platform/go-auth-api/internal/models"
)

type Queries struct{
	Existing string
	Insert string			
	GetUser string	
}

func GetQueries() (*Queries){
	return &Queries{Existing: `
	SELECT COUNT(*) FROM users WHERE email = $1 LIMIT 1
	`,
	Insert:`
		INSERT INTO users (user_id, username, email) 
		VALUES ($1, $2, $3);
	`,
	GetUser: `SELECT * FROM users WHERE email = $1 LIMIT 1`,
}
}

func CheckExisting(email string, ctx context.Context,db *sql.DB)(bool, error){
	var count int
	err := db.QueryRowContext(ctx,GetQueries().Existing,email).Scan()
	if err != nil{
		return false, err
	}
	return count >0, err

}

func GetUsername(email string,ctx context.Context,db *sql.DB) (string, error){
	var user models.Users
	err := db.QueryRowContext(ctx, GetQueries().GetUser, email).Scan(&user.Id,&user.Username,&user.Email)
	if err != nil{
		return "",err
	}
	return user.Username, nil
}
