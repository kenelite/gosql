package storage

import (
	"fmt"
	"sync"
)

type Storage struct {
	tables map[string]*Table
	mu     sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		tables: make(map[string]*Table),
	}
}

func (s *Storage) CreateTable(name string, cols []Column) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tables[name]; exists {
		return fmt.Errorf("table %s already exists", name)
	}
	s.tables[name] = NewTable(name, cols)
	return nil
}

func (s *Storage) Insert(table string, row Row) error {
	s.mu.RLock()
	t, ok := s.tables[table]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("table %s not found", table)
	}
	return t.Insert(row)
}

func (s *Storage) SelectAll(table string) ([]Column, []Row, error) {
	s.mu.RLock()
	t, ok := s.tables[table]
	s.mu.RUnlock()
	if !ok {
		return nil, nil, fmt.Errorf("table %s not found", table)
	}
	return t.Columns, t.SelectAll(), nil
}
