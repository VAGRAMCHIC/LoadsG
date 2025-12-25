package model


//----------- User ------------------

type User struct {
	UID       string `json:"uid" binding:"required"`
	PasswordHash string `json:"password" binding:"required"`
}


var DB_TABLE_USERS = DB_TABLE{
	Name: "users",
	Params: map[string]string{
		"id":       		 "SERIAL PRIMARY KEY",
		"uid":					 "TEXT NOT NULL",
		"password_hash": "TEXT NOT NULL",
	},
}


