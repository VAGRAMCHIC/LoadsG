package lib

import (
	"context"
	"log"
	"fmt"
	"sort"
	"strings"
	"github.com/jackc/pgx/v5"
)

func InitDB(conn *pgx.Conn) (bool, error) {
	code, err := initTable(conn, DB_TABLE_USERS)
	if err != nil {
		return code, err
	}
	return code, err
}

func initTable(conn *pgx.Conn, table DB_TABLE) (bool, error) {
	var init_table_line string

	var keys []string
	for k := range table.Params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", table.Name))

	for i, k := range keys {
		b.WriteString(fmt.Sprintf("    %s %s", k, table.Params[k]))
		if i < len(keys)-1 {
			b.WriteString(",\n") // запятая только между полями
		} else {
			b.WriteString("\n") // без запятой на последней строке
		}
	}
	b.WriteString(");")

	init_table_line = b.String()
	fmt.Printf("Create table %s: %s", table.Name, init_table_line)
	_, err := conn.Exec(context.Background(), init_table_line)
	//defer conn.Close(context.Background())
	if err != nil {
		log.Printf("cant create table: %s", err.Error())
		return false, err
	}
	return true, err
}



func Connect(pgConn string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), pgConn)
	if err != nil {
		log.Printf("cant create connection: %s", err.Error())
	}
	return conn
}

func InsertUser(conn *pgx.Conn, user User) {
	_, err := conn.Exec(context.Background(),
		"INSERT INTO users (username, password) VALUES ($1, $2)", user.Id, user.Password)
	if err != nil {
		log.Fatalf("Insert data error: %s", err.Error())
		defer conn.Close(context.Background())
	}
	defer conn.Close(context.Background())
}

func GetUser(conn *pgx.Conn, username string) (User, error) {
	var user User
	err := conn.QueryRow(context.Background(), "SELECT * FROM users where username=$1", username).Scan(&user.Id, &user.Password)
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
