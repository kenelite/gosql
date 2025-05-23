package storage

import (
	"fmt"
	"sync"
)

type Table struct {
	Name    string
	Columns []Column
	Rows    []Row
	mu      sync.RWMutex
}

func NewTable(name string, cols []Column) *Table {
	return &Table{
		Name:    name,
		Columns: cols,
		Rows:    make([]Row, 0),
	}
}

func (t *Table) Insert(row Row) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(row) != len(t.Columns) {
		return fmt.Errorf("column count mismatch")
	}
	t.Rows = append(t.Rows, row)
	return nil
}

func (t *Table) SelectAll() []Row {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return append([]Row(nil), t.Rows...) // copy
}
