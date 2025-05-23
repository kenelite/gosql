package storage

type Column struct {
	Name string
	Type string // e.g., "INT", "VARCHAR"
}

type Row []interface{}
