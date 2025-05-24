package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const storageDir = "data"

type FileStore struct {
	tables map[string]*Table
	mu     sync.RWMutex
}

func NewFileStore(storageDir string) (*FileStore, error) {
	fs := &FileStore{
		tables: make(map[string]*Table),
	}

	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}

	return fs, fs.loadTables()
}

func (fs *FileStore) loadTables() error {
	entries, err := os.ReadDir(storageDir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		path := filepath.Join(storageDir, e.Name())
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		var t Table
		if err := json.NewDecoder(file).Decode(&t); err != nil {
			return err
		}
		t.mu = sync.RWMutex{}
		fs.tables[t.Name] = &t
	}
	return nil
}

func (fs *FileStore) saveTable(t *Table) error {
	path := filepath.Join(storageDir, t.Name+".json")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	t.mu.RLock()
	defer t.mu.RUnlock()

	return json.NewEncoder(file).Encode(t)
}

func (fs *FileStore) CreateTable(name string, cols []Column) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, exists := fs.tables[name]; exists {
		return fmt.Errorf("table %s already exists", name)
	}
	t := NewTable(name, cols)
	fs.tables[name] = t
	return fs.saveTable(t)
}

func (fs *FileStore) Insert(table string, row Row) error {
	fs.mu.RLock()
	t, ok := fs.tables[table]
	fs.mu.RUnlock()
	if !ok {
		return fmt.Errorf("table %s not found", table)
	}
	if err := t.Insert(row); err != nil {
		return err
	}
	return fs.saveTable(t)
}

func (fs *FileStore) SelectAll(table string) ([]Column, []Row, error) {
	fs.mu.RLock()
	t, ok := fs.tables[table]
	fs.mu.RUnlock()
	if !ok {
		return nil, nil, fmt.Errorf("table %s not found", table)
	}
	return t.Columns, t.SelectAll(), nil
}
