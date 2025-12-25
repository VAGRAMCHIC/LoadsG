package postgres

import (
	"context"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)


type UserRepository struct{
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &UserRepository{
		db:db,
	}
}

func (r *UserRepository) GetById(ctx context.Context, uid string) (*model.User, error){
	var user model.User
	err := r.db.QueryRow(context.Background(), "SELECT * FROM users where uid=$1", uid).Scan(&user.UID, &user.PasswordHash)
	if err != nil{
		return &user, err
	}
	return &user, nil
	
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error){
	_, err := r.db.Exec(context.Background(),
		"INSERT INTO users (uid, password) VALUES ($1, $2)", user.UID, user.PasswordHash)
	if err != nil {
		log.Fatalf("Insert data error: %s", err.Error())
	}
	return user, err
}


