package model


type User struct {
	UID       string `json:"uid" binding:"required"`
	Token 		string `json:"token" binding:"required"`
	Comment 	string `json:"comment"`
}


var DB_TABLE_USERS = DB_TABLE{
	Name: "users",
	Params: map[string]string{
		"uid":       		 "UUID PRIMARY KEY DEFAULT gen_random_uuid()",
		"token_hash": "TEXT NOT NULL",
		"comment": "TEXT",
	},
	InlinePrimaryKey: true,
}


