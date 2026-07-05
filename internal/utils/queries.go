package utils

import (
	"context"
	"database/sql"
)

type Queries struct{
	Existing string
	Insert string				
}

func GetQueries() (*Queries){
	return &Queries{Existing: `
	SELECT COUNT(*) FROM users WHERE email = $1 LIMIT 1
	`,
	Insert:`
		INSERT INTO users (user_id, username, email) 
		VALUES ($1, $2, $3);
	`,
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
