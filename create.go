package main

import (
	"flag"
	"log"
	"os"
	"path"
	"time"

	"github.com/nbd-wtf/go-nostr"
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
func (s *Create) Run(cfg *Config) error {

	if s.title == "" {
		return ErrNoTitle
	}
	if s.content == "" {
		return ErrNoContent
	}

	e := s.createEvent()

	err := publishEvent(e, cfg.Nsec, cfg.Relays)
	if err != nil {
		return err
	}

	return nil
}

func (s Create) createEvent() *nostr.Event {

	var tags nostr.Tags

	identifier := time.Now().Format("200601021504")
	tags = append(tags, nostr.Tag{"d", identifier})
	tags = append(tags, nostr.Tag{"title", s.title})

	for _, v := range s.hashtags {
		tags = append(tags, nostr.Tag{"t", v})
	}

	return &nostr.Event{
		Kind:      nostr.KindArticle,
		Content:   s.content,
		CreatedAt: nostr.Now(),
		Tags:      tags,
	}
}
