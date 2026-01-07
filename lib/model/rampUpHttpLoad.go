package model

type RampUpHttpLoad struct {
	Id        string            `json:"id" binding:"required"`
	LoadJobId string            `json:"load_job_id" binding:"required"`
	RPS_S     int               `json:"rps_s" binding:"required"`
	RPS_F     int               `json:"rps_f" binding:"required"`
	Duration  float32           `json:"duration" binding:"required"`
	Payload   map[string]string `json:"payload" binding:"required"`
}

var DB_TABLE_RAMPUP_HTTP_LOAD = DB_TABLE{
	Name: "ramp_up_http_load",
	Params: map[string]string{
		"load_job_id": "UUID NOT NULL",
		"duration":    "REAL NOT NULL",
		"rps_s":       "INTEGER NOT NULL",
		"rps_f":       "INTEGER NOT NULL",
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
