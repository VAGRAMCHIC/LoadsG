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
	err := r.db.QueryRow(context.Background(), "SELECT id, iss, uid, token_hash, created_at, expires_at, revoked FROM users where token_hash=$1", hash).Scan(&token.Id, &token.Uid, &token.TokenHash,
		&token.Issuer, &token.Created_at, &token.Expires_at, &token.Revoked)
	if err != nil{
		log.Printf("cant get refresh token by hash: %s", err)
		return &token, err
	}
	return &token, nil
	
}

func (r *JWTRefreshRepository) Create(ctx context.Context, jwt *model.JWTRefreshToken) (*model.JWTRefreshToken, error){
	_, err := r.db.Exec(context.Background(),
		"INSERT INTO refresh_token (uid, token_hash, iss, created_at, expires_at, revoked) VALUES ($1, $2, $3, $4, $5, $6)", jwt.Uid, jwt.TokenHash,
			jwt.Issuer, jwt.Created_at, jwt.Expires_at, jwt.Revoked)
	if err != nil {
		log.Fatalf("Insert data error: %s", err.Error())
	}
	return jwt, err
}


