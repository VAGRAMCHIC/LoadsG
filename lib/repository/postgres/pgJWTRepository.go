package postgres


import (
	"context"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"
	"time"
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
	err := r.db.QueryRow(context.Background(), "SELECT id, user_uid, token_hash, expires_at FROM refresh_tokens where token_hash=$1", hash).Scan(&token.Id, &token.UserUid, &token.TokenHash, &token.ExpiresAt)
	if err != nil{
		log.Printf("cant get refresh token by hash: %s", err)
		return &token, err
	}
	return &token, nil
	
}

func (r *JWTRefreshRepository) Create(ctx context.Context, jwt *model.JWTRefreshToken) (*model.JWTRefreshToken, error){
	err := r.db.QueryRow(context.Background(),
		"INSERT INTO refresh_tokens (user_uid, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING user_uid, token_hash, expires_at", jwt.UserUid, jwt.TokenHash, jwt.ExpiresAt).Scan(&jwt.UserUid, &jwt.TokenHash, &jwt.ExpiresAt)
	if err != nil {
		log.Fatalf("Insert data error: %s", err)
	}
	return jwt, nil
}



func (r *JWTRefreshRepository) Save(ctx context.Context, userID string,
	hash string, expires time.Time,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (user_uid, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, hash, expires)
	return err
}

func (r *JWTRefreshRepository) Get(ctx context.Context, hash string) (string, error) {
	var uid string
	err := r.db.QueryRow(ctx, `
		SELECT user_uid FROM refresh_tokens
		WHERE token_hash = $1
	`, hash).Scan(&uid)
	return uid, err
}




func (r *JWTRefreshRepository) Delete(ctx context.Context, jwtHash string) error{
	cmd, err := r.db.Exec(context.Background(),
		"DELETE FROM refresh_tokens WHERE token_hash=$1", jwtHash)
	if err != nil {
		log.Fatalf("Delete data error: %s", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		return nil
	}

	return nil
}
