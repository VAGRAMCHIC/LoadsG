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

func (r *UserRepository) GetByUID(ctx context.Context, uid string) (*model.User, error){
	var user model.User
	err := r.db.QueryRow(context.Background(), "SELECT uid, token_hash, comment FROM users where uid=$1", uid).Scan(&user.UID, &user.Token, &user.Comment)
	if err != nil{
		log.Printf("cant get user by uid: %s", err)
		return &user, err
	}
	return &user, nil
	
}


func (r *UserRepository) GetByComment(ctx context.Context, comment string) (*model.User, error){
	var user model.User
	err := r.db.QueryRow(context.Background(), "SELECT uid, token_hash, comment FROM users where comment=$1", comment).Scan(&user.UID, &user.Token, &user.Comment)
	if err != nil{
		log.Printf("cant get user by uid: %s", err)
		return &user, err
	}
	return &user, nil
	
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error){
	if user.Comment == "" {
		err := r.db.QueryRow(context.Background(),
			"INSERT INTO users (token_hash) VALUES ($1) RETURNING uid", user.Token).Scan(&user.UID)
		if err != nil {
			log.Fatalf("Insert data error: %s", err.Error())
			return user, err
		}
		return user, nil
	}
	
	err := r.db.QueryRow(context.Background(),
		"INSERT INTO users (comment, token_hash) VALUES ($1, $2) RETURNING uid, comment", user.Comment, user.Token).Scan(&user.UID, &user.Comment)
	log.Print(user)
	if err != nil {
		log.Fatalf("Insert data error: %s", err.Error())
		return user, err
	}
	
	return user, nil
}


