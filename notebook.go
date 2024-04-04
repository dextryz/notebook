package main

import (
	"fmt"
	"log"
	"os"
)

// 1. Notebook abstracts nostr, and therefore does not work with nostr.Event, rather with nz.Note

// Abstract working with nostr relays and local embedded database
type Notebook struct {
	cfg *Config
	dir string
}

func NewNotebook(dir string) *Notebook {

	// TODO: Handle the case for when dir already exists
	err := os.Mkdir(dir, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	return &Notebook{
		dir: dir,
	}
}

// 1. Store event in local working directory
// 2. Store event in SQLite for fast querying
func (s Notebook) store(n Note) error {
	filename := fmt.Sprintf("%s/%s.md", s.dir, n.Identifier())
	err := os.WriteFile(filename, []byte(n.Content), 0660)
	if err != nil {
		return err
	}
	return nil
}
