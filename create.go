package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

// Implemenets the Runner interface

type Create struct {
	fs *flag.FlagSet

	title    string
	content  string
	hashtags []string
}

func NewCreate() *Create {

	s := Create{
		fs: flag.NewFlagSet("new", flag.ExitOnError),
	}

	s.fs.StringVar(&s.title, "title", "", "the filename to be committed")

	s.fs.Func("tag", "tag article", func(v string) error {
		s.hashtags = append(s.hashtags, v)
		return nil
	})

	return &s
}

func (s *Create) Parse(args []string) error {

	// TODO Improve
	filename := args[0]
	if path.Ext(filename) != ".md" {
		log.Fatalln("file type has to be .md")
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	s.content = string(data)

	return s.fs.Parse(args[1:])
}

func (s *Create) Name() string {
	return s.fs.Name()
}

// 1. Check that title and content is non-empty
// 2. Create the nostr event from the commandline arguments
// 3. Publish the event to the config relays
func (s *Create) Run(container *Container) error {

	// TODO move this to note opts
	if s.title == "" {
		return ErrNoTitle
	}
	if s.content == "" {
		return ErrNoContent
	}

	// Get the current notebook instance
	nb, err := container.CurrentNotebook()
	if err != nil {
		return err
	}
	if nb == nil {
		log.Fatalln("no notebook specified")
	}

	// Add note to this notebook and publish to nostr relays
	note, err := nb.Add(s.title, s.content, s.hashtags)
	if err != nil {
		return err
	}

	fmt.Printf("Note created, open with editor: %s", note.Path)

	return nil
}
