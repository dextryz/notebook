package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
)

// Usage:
//
// > nz list

type List struct {
	fs *flag.FlagSet
}

func NewList() *List {

	s := List{
		fs: flag.NewFlagSet("list", flag.ExitOnError),
	}

	return &s
}

func (s *List) Parse(args []string) error {
	return s.fs.Parse(args)
}

func (s *List) Name() string {
	return s.fs.Name()
}

// 1. Check that title and content is non-empty
// 2. Create the nostr event from the commandline arguments
// 3. Publish the event to the config relays
func (s *List) Run(cfg *Config) error {

	ctx := context.Background()

	events, err := requestSortedEvents(ctx, cfg.Nsec, cfg.Relays)
	if err != nil {
		return err
	}

	ViewTitle(events)

	return nil
}

func ViewTitle(events []*nostr.Event) {

	for _, e := range events {

		title, identifier := "", ""

		for _, t := range e.Tags {
			if t.Key() == "title" {
				title = t.Value()
			}
			if t.Key() == "d" {
				identifier = t.Value()
			}
		}

		if identifier == "" {
			fmt.Printf("Article with title (%s) has no identifier\n", title)
		} else {
			fmt.Printf("%s\n", identifier)
			fmt.Printf("- %s\n", title)
			fmt.Println("")
		}
	}
}
