package model

type ConstantHttpLoad struct {
	Id        string            `json:"id" binding:"required"`
	LoadJobId string            `json:"load_job_id" binding:"required"`
	Count     int               `json:"count" binding:"required"`
	Payload   map[string]string `json:"payload" binding:"required"`
}

var DB_TABLE_CONSTANT_HTTP_LOAD = DB_TABLE{
	Name: "constant_http_load",
	Params: map[string]string{
		"load_job_id": "UUID NOT NULL",
		"count":       "INTEGER NOT NULL",
		"payload":     "JSONB NOT NULL",
	},
	PrimaryKey: []string{"load_job_id"},
	ForeignKeys: []ForeignKey{
		{
			Column:    "load_job_id",
			RefTable:  "load_job",
			RefColumn: "id",
			OnDelete:  "CASCADE",
		},
	},
}
