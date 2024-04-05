package main

import (
	"log"
	"os"
)

type Container struct {
	cfg      *Config
	notebook *Notebook
}

func NewContainer(cfg *Config) *Container {
	return &Container{
		cfg:      cfg,
		notebook: nil,
	}
}

func (s *Container) CurrentNotebook() (*Notebook, error) {

	name := os.Getenv("NOTEBOOK")
	dir := os.Getenv("NOTEBOOK_DIR")

	if name != "" && dir != "" {
		s.notebook = NewNotebook(s.cfg, name, dir)
		log.Printf("current notebook [%s] found with dir: %s", name, dir)
	}

	return s.notebook, nil
}
