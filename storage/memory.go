package storage

type Table struct {
	Columns []string
	Rows    [][]interface{}
}

var db = map[string]*Table{}
