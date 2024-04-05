package main

import (
	"flag"
	"log"
	"os"
)

// Creates a notebook structure that abstracts the
// local working directory.

type Init struct {
	fs *flag.FlagSet

	name string
	dir  string
}

func NewInit() *Init {

	s := Init{
		fs: flag.NewFlagSet("init", flag.ExitOnError),
	}

	s.fs.StringVar(&s.name, "name", "", "the filename to be committed")
	s.fs.StringVar(&s.dir, "dir", "", "the filename to be committed")

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

	os.Setenv("NOTEBOOK", s.name)
	os.Setenv("NOTEBOOK_DIR", s.dir)

	nb := NewNotebook(container.cfg, s.name, s.dir)

	notes, err := nb.Init()
	if err != nil {
		return err
	}

	// Add the newly created notebook to the container
	container.notebook = nb

	log.Printf("notebook initiated with %d notes pulled from relays", len(notes))

	return nil
}
