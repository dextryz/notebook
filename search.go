package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/nbd-wtf/go-nostr"
)

// Implemenets the Runner interface

type Search struct {
	fs *flag.FlagSet

	title    string
	hashtags []string
}

func NewSearch() *Search {

	s := Search{
		fs: flag.NewFlagSet("search", flag.ExitOnError),
	}

	s.fs.StringVar(&s.title, "title", "", "article title")

	s.fs.Func("tag", "tag article", func(v string) error {
		s.hashtags = append(s.hashtags, v)
		return nil
	})

	return &s
}

func (s *Search) Parse(args []string) error {
	return s.fs.Parse(args)
}

func (s *Search) Name() string {
	return s.fs.Name()
}

// 1. Check that title and content is non-empty
// 2. Create the nostr event from the commandline arguments
// 3. Publish the event to the config relays
func (s *Search) Run(container *Container) error {

	nb, err := container.CurrentNotebook()
	if err != nil {
		return err
	}

	tags := make(nostr.TagMap)
	if s.title != "" {
		tags["title"] = []string{s.title}
	}
	tags["t"] = append(tags["t"], s.hashtags...)

	filter := nostr.Filter{
		Kinds: []int{nostr.KindArticle},
		Tags:  tags,
		Limit: 500,
	}

	notes, err := nb.Search(filter)
	if err != nil {
		return err
	}

	slog.Info("events found", "eventCount", len(notes))

	fmt.Println("\n[*] Notes")
	for _, n := range notes {
		fmt.Printf("%s\n", n.Title())
		fmt.Printf("%s\n\n", n.Path)
	}

	return nil
}
