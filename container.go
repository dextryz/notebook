package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fiatjaf/eventstore"
)

type Container struct {
	cfg      *Config
	db       eventstore.Store
	notebook *Notebook
}

func NewContainer(cfg *Config, db eventstore.Store) *Container {
	return &Container{
		cfg:      cfg,
		notebook: nil,
		db:       db,
	}
}

func (s *Container) CurrentNotebook() (*Notebook, error) {

	dir, ok := os.LookupEnv("NOTEBOOK")
	if !ok {
		return nil, fmt.Errorf("NOTEBOOK env var not set")
	}

	s.notebook = NewNotebook(s.cfg, s.db, dir)
	log.Printf("current notebook found with path: %s", dir)

	return s.notebook, nil
}
