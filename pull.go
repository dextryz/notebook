package main

import (
	"context"
	"flag"
	"log"
)

// Creates a notebook structure that abstracts the
// local working directory.

type Pull struct {
	fs *flag.FlagSet

	dir string
}

func NewPull() *Pull {

	s := Pull{
		fs: flag.NewFlagSet("pull", flag.ExitOnError),
	}

	return &s
}

func (s *Pull) Parse(args []string) error {

	// TODO move to newPull using flag.Args()
	s.dir = args[0]

	return s.fs.Parse(args[1:])
}

func (s *Pull) Name() string {
	return s.fs.Name()
}

// 1. Creates a new Notebook
// 2. Pulls all the author articles from nostr relays
// 3. Stores each event in the notebook as a note
func (s *Pull) Run(cfg *Config) error {

	ctx := context.Background()

	events, err := requestSortedEvents(ctx, cfg.Nsec, cfg.Relays)
	if err != nil {
		return err
	}

	log.Printf("%d events pulled from relays", len(events))

	nb := NewNotebook(s.dir)
	for _, v := range events {
		n := Note{
			Event: v,
		}
		err := nb.store(n)
		if err != nil {
			return err
		}
	}

	return nil
}
