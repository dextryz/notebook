package main

import (
	"flag"
	"fmt"
)

// Usage:
//
// > nz list --format short

type List struct {
	fs *flag.FlagSet

	format string // short, long, json
	header string // notebook details, metadata
	footer string // note count, metadata, ----
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
func (s *List) Run(container *Container) error {

	// Get the current notebook instance
	nb, err := container.CurrentNotebook()
	if err != nil {
		return err
	}

	notes, err := nb.FindNotes()

	for _, n := range notes {
		format(n)
	}

	return nil
}

// TODO Make this a template and styler pattern
func format(n *Note) {
	fmt.Printf("%s\n", n.Identifier())
	fmt.Printf("- %s\n", n.Title())
	fmt.Println("")
}
