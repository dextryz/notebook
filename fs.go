package main

// Context
//
// Only stores, deletes, and updates the physical markdown files.
// We still need a caching layer to store note metadata

import (
	"fmt"
	"os"
)

// TODO
type FileType int

type FileStorage struct {
	dir string
}

// TODO: Do I have to init, or can I just mkdir when I write?
func (s *FileStorage) Init() error {

	// TODO: Handle the case for when dir already exists
	err := os.Mkdir(s.dir, 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

// Read file content from local filesystem
func (s *FileStorage) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// 1. Store event in local working directory
// 2. Store event in SQLite for fast querying
func (s *FileStorage) Store(n *Note) error {

	n.Path = fmt.Sprintf("%s/%s.md", s.dir, n.Identifier())

	err := os.WriteFile(n.Path, []byte(n.Content), 0660)
	if err != nil {
		return err
	}

	return nil
}
