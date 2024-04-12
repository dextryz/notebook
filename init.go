package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

// Creates a notebook structure that abstracts the local working directory.

type Init struct {
	fs *flag.FlagSet
}

func NewInit() *Init {

	s := Init{
		fs: flag.NewFlagSet("init", flag.ExitOnError),
	}

	return &s
}

func (s *Init) Parse(args []string) error {
	return s.fs.Parse(args)
}

func (s *Init) Name() string {
	return s.fs.Name()
}

// 1. Creates a new Notebook
// 2. Inits all the author articles from nostr relays
// 3. Stores each event in the notebook as a note
func (s *Init) Run(container *Container) error {

	dir, ok := os.LookupEnv("NOTEBOOK")
	if !ok {
		return fmt.Errorf("NOTEBOOK env var not set")
	}

	nb := NewNotebook(container.cfg, container.db, dir)

	notes, err := nb.Init()
	if err != nil {
		return err
	}

	// Add the newly created notebook to the container
	container.notebook = nb

	slog.Info("notebook initiated with %d notes pulled from relays", "noteCount", len(notes))

	return nil
}
