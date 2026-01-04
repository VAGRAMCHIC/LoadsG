package postgres


import (
	"context"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)


type JWTRefreshRepository struct{
	db *pgxpool.Pool
}

func NewJWTRefreshRepository(db *pgxpool.Pool) repository.JWTRefreshRepository {
	return &JWTRefreshRepository{
		db:db,
	}
}

func (r *JWTRefreshRepository) GetByJWTHash(ctx context.Context, hash string) (*model.JWTRefreshToken, error){
	var token model.JWTRefreshToken
	err := r.db.QueryRow(context.Background(), "SELECT id, user_uid, token_hash FROM refresh_tokens where token_hash=$1", hash).Scan(&token.Id, &token.UserUid, &token.TokenHash)
	if err != nil{
		log.Printf("cant get refresh token by hash: %s", err)
		return &token, err
	}
	return &token, nil
	
}

func (r *JWTRefreshRepository) Create(ctx context.Context, jwt *model.JWTRefreshToken) (*model.JWTRefreshToken, error){
	_, err := r.db.Exec(context.Background(),
		"INSERT INTO refresh_tokens (uid, token_hash) VALUES ($1, $2)", jwt.UserUid, jwt.TokenHash)
	if err != nil {
		log.Fatalf("Insert data error: %s", err.Error())
	}
	return jwt, err
}
