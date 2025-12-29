package model

import "time"

type JWTRefreshToken struct {
	Id 				int `json:"id" binding:"required"`
	Uid				string `json:"uid" binding:"required"`
	TokenHash string `json:"jwt_hash" binding:"required"`
	Issuer 		string `json:"iss" binding:"required"`
	Created_at time.Time `json:"created_at" binding:"required"`
	Expires_at time.Time `json:"expires_at" binding:"required"`
	Revoked 	bool `json:"revoked" binding:"required"`
		
}

var DB_TABLE_JWTTOKENS = DB_TABLE{
	Name: "refresh_tokens",
	Params: map[string]string{
		"id": "UUID NOT NULL",
		"iss": "TEXT NOT NULL",
		"uid": "UUID NOT NULL",
		"token_hash": "TEXT NOT NULL",
		"created_at" : "TIMESTAMP DEFAULT NOW()",
		"expires_at" : "TIMESTAMP NOT NULL",
		"revoked" : "BOOLEAN DEFAULT FALSE",
	},
	PrimaryKey: []string{"id"},
	ForeignKeys: []ForeignKey{
		{
			Column:    "uid",
			RefTable:  "users",
			RefColumn: "id",
			OnDelete:  "CASCADE",
		},
	},
}
