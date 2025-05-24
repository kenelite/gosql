package parser

import (
	"errors"
	"fmt"
	"github.com/kenelite/gosql/storage"
	"regexp"
	"strings"
)

type Statement interface {
	Type() string
}

type CreateTableStmt struct {
	TableName string
	Columns   []storage.Column
}

func (s *CreateTableStmt) Type() string { return "CREATE" }

type InsertStmt struct {
	TableName string
	Values    []interface{}
}

func (s *InsertStmt) Type() string { return "INSERT" }

type SelectStmt struct {
	TableName string
}

func (s *SelectStmt) Type() string { return "SELECT" }

// Parse parses a SQL query into a Statement object
func Parse(query string) (Statement, error) {
	query = strings.TrimSpace(query)
	lower := strings.ToLower(query)

	switch {
	case strings.HasPrefix(lower, "create table"):
		return parseCreateTable(query)
	case strings.HasPrefix(lower, "insert into"):
		return parseInsert(query)
	case strings.HasPrefix(lower, "select"):
		return parseSelect(query)
	default:
		return nil, errors.New("unsupported statement")
	}
}

func parseCreateTable(query string) (Statement, error) {
	// Example: CREATE TABLE users (id INT, name VARCHAR)
	re := regexp.MustCompile(`(?i)^create table (\w+) *\((.+)\)`)
	matches := re.FindStringSubmatch(query)
	if len(matches) != 3 {
		return nil, errors.New("invalid CREATE TABLE syntax")
	}

	tableName := matches[1]
	colsDef := matches[2]
	cols := []storage.Column{}
	for _, col := range strings.Split(colsDef, ",") {
		parts := strings.Fields(strings.TrimSpace(col))
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid column definition: %s", col)
		}
		cols = append(cols, storage.Column{Name: parts[0], Type: strings.ToUpper(parts[1])})
	}

	return &CreateTableStmt{TableName: tableName, Columns: cols}, nil
}

func parseInsert(query string) (Statement, error) {
	// Example: INSERT INTO users VALUES (1, 'Alice')
	re := regexp.MustCompile(`(?i)^insert into (\w+) values *\((.+)\)`)
	matches := re.FindStringSubmatch(query)
	if len(matches) != 3 {
		return nil, errors.New("invalid INSERT syntax")
	}
	table := matches[1]
	valsRaw := matches[2]
	valTokens := strings.Split(valsRaw, ",")
	row := make([]interface{}, 0, len(valTokens))
	for _, tok := range valTokens {
		tok = strings.TrimSpace(tok)
		if strings.HasPrefix(tok, "'") && strings.HasSuffix(tok, "'") {
			row = append(row, strings.Trim(tok, "'"))
		} else {
			row = append(row, parseInt(tok))
		}
	}
	return &InsertStmt{TableName: table, Values: row}, nil
}

func parseSelect(query string) (Statement, error) {
	// Example: SELECT * FROM users
	tokens := strings.Fields(strings.ToLower(query))
	if len(tokens) < 4 || tokens[1] != "*" || tokens[2] != "from" {
		return nil, errors.New("invalid SELECT syntax")
	}
	table := tokens[3]
	return &SelectStmt{TableName: table}, nil
}

func parseInt(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
