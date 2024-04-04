package main

import (
	"flag"
	"log"
)

// Creates a notebook structure that abstracts the
// local working directory.

type Init struct {
	fs *flag.FlagSet

	dir string
}

func NewInit() *Init {
	s := Init{
		fs: flag.NewFlagSet("pull", flag.ExitOnError),
	}
	return &s
}

// 1. Parse the arguments then parse the flags
func (s *Init) Parse(args []string) error {
	s.dir = args[0]
	return s.fs.Parse(args[1:])
}

func (s *Init) Name() string {
	return s.fs.Name()
}

// 1. Creates a new Notebook
// 2. Inits all the author articles from nostr relays
// 3. Stores each event in the notebook as a note
func (s *Init) Run(container *Container) error {

	nb := NewNotebook(container.cfg, s.dir)

	notes, err := nb.Init()
	if err != nil {
		return err
	}

	log.Printf("notebook initiated with %d notes pulled from relays", len(notes))

	// Add the newly created notebook to the container
	container.notebook = nb

	return nil
}
