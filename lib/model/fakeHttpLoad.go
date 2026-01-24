package model

type FakeHttpLoad struct {
	Id        string  `json:"id" binding:"required"`
	LoadJobId string  `json:"load_job_id" binding:"required"`
	Duration  float32 `json:"duration" binding:"required"`
}

var DB_TABLE_FAKE_HTTP_LOAD = DB_TABLE{
	Name: "fake_http_load",
	Params: map[string]string{
		"load_job_id": "UUID NOT NULL",
		"duration":    "REAL NOT NULL",
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
