package model

type FixedHttpLoad struct {
	Id 						string 						`json:"id" binding:"required"`
	LoadJobId 		string						`json:"load_job_id" binding:"required"`
	RequestCount	int 							`json:"count" binding:"required"`
	Payload				map[string]string `json:"payload" binding:"required"`
	
}


var DB_TABLE_FIXED_HTTP_LOAD = DB_TABLE{
	Name: "fixed_http_load",
	Params: map[string]string{
		"id": 						"SERIAL PRIMARY KEY",
		"load_job_id":		"UUID NOT NULL",
		"request_count":	"INTEGER NOT NULL",
		"payload": 				"JSONB NOT NULL",
	},
	InlinePrimaryKey: true,
	ForeignKeys: []ForeignKey{
		{
			Column: "_id",
			RefTable:"load_job",
			RefColumn: "id",
			OnDelete: "CASCADE",
		},
	},
}
