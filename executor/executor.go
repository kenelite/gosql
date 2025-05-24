package executor

import (
	"errors"
	"github.com/kenelite/gosql/parser"
	"github.com/kenelite/gosql/protocol"
	"github.com/kenelite/gosql/storage"
)

type Executor struct {
	Store *storage.FileStore
}

func NewExecutor(store *storage.FileStore) *Executor {
	return &Executor{Store: store}
}

func (e *Executor) Execute(stmt parser.Statement, conn *protocol.Conn) error {
	switch s := stmt.(type) {
	case *parser.CreateTableStmt:
		return e.execCreateTable(s, conn)
	case *parser.InsertStmt:
		return e.execInsert(s, conn)
	case *parser.SelectStmt:
		return e.execSelect(s, conn)
	default:
		return errors.New("unsupported statement type")
	}
}

func (e *Executor) execCreateTable(stmt *parser.CreateTableStmt, conn *protocol.Conn) error {
	err := e.Store.CreateTable(stmt.TableName, stmt.Columns)
	if err != nil {
		return conn.WriteError(1064, err.Error())
	}
	return conn.WriteOK()
}

func (e *Executor) execInsert(stmt *parser.InsertStmt, conn *protocol.Conn) error {
	err := e.Store.Insert(stmt.TableName, stmt.Values)
	if err != nil {
		return conn.WriteError(1064, err.Error())
	}
	return conn.WriteOK()
}

func (e *Executor) execSelect(stmt *parser.SelectStmt, conn *protocol.Conn) error {
	cols, rows, err := e.Store.SelectAll(stmt.TableName)
	if err != nil {
		return conn.WriteError(1064, err.Error())
	}
	names := make([]string, len(cols))
	for i, col := range cols {
		names[i] = col.Name
	}
	return conn.WriteResultSet(names, rows)
}
