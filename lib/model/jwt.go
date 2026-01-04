package model


type JWTRefreshToken struct {
	Id 				int `json:"id" binding:"required"`
	UserUid				string `json:"uid" binding:"required"`
	TokenHash string `json:"jwt_hash" binding:"required"`
}

var DB_TABLE_JWTTOKENS = DB_TABLE{
	Name: "refresh_tokens",
	Params: map[string]string{
		"id":		      "UUID PRIMARY KEY DEFAULT gen_random_uuid()",
		"user_uid": 				"UUID NOT NULL",
		"token_hash": "TEXT NOT NULL",
	},
	InlinePrimaryKey: true,
	ForeignKeys: []ForeignKey{
		{
			Column:    "user_uid",
			RefTable:  "users",
			RefColumn: "uid",
			OnDelete:  "CASCADE",
		},
	},
}
