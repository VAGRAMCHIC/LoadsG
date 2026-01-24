package model

type Event struct {
	Id        string `json:"id" binding:"required"`
	Status    string `json:"status" binding:"required"`
	LoadJobId string `json:"load_job_id" binding:"required"`
}

var DB_TABLE_EVENTS = DB_TABLE{
	Name: "events",
	Params: map[string]string{
		"status":      "TEXT NOT NULL DEFAULT 'pending'",
		"load_job_id": "UUID NOT NULL",
	},
	InlinePrimaryKey: true,
	ForeignKeys: []ForeignKey{
		{
			Column:    "load_job_id",
			RefTable:  "load_job",
			RefColumn: "id",
			OnDelete:  "CASCADE",
		},
	},
}
