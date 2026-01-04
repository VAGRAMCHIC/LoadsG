package model


type ForeignKey struct {
	Column    string
	RefTable  string
	RefColumn string
	OnDelete  string // CASCADE | SET NULL | RESTRICT
}


type DB_TABLE struct {
	Name   string            `json:"name" binding:"required"`
	Params map[string]string `json:"params" binding:"required"`
	PrimaryKey  []string
	InlinePrimaryKey bool
	ForeignKeys []ForeignKey
}

