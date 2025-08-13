package lib

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func InitDB(conn pgx.Conn) bool {
	_, err := conn.Exec(context.Background(), `
			CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				username TEXT NOT NULL,
				password TEXT NOT NULL,
			)
		`)
	if err != nil {
		log.Printf("cant init database: %s", err.Error())
		defer conn.Close(context.Background())
		return false
	}
	defer conn.Close(context.Background())
	return true
}

func connect(pgConn string) pgx.Conn {
	conn, err := pgx.Connect(context.Background(), pgConn)
	if err != nil {
		log.Printf("cant create connection: %s", err.Error())
	}
	return *conn
}

func InsertUser(conn pgx.Conn, user User) {
	_, err := conn.Exec(context.Background(),
		"INSERT INTO users (username, password) VALUES ($1, $2)", user.UserName, user.Password)
	if err != nil {
		log.Fatalf("Insert data error: %s", err.Error())
		defer conn.Close(context.Background())
	}
	defer conn.Close(context.Background())
}

func GetUser(conn pgx.Conn, username string) (User, error) {
	var user User
	err := conn.QueryRow(context.Background(), "SELECT * FROM users where username=$1", username).Scan(&user.UserName, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Пользователь не найден
			return user, nil
		}
		return user, err
	}
	defer conn.Close(context.Background())
	return user, err
}
