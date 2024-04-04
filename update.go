package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

// > nz update /tmp/zk/202404041455.md

type Update struct {
	fs *flag.FlagSet

	filename string
}

func NewUpdate() *Update {

	s := Update{
		fs: flag.NewFlagSet("update", flag.ExitOnError),
	}

	return &s
}

func (s *Update) Parse(args []string) error {

	s.filename = args[1]

	return s.fs.Parse(args[1:])
}

func (s *Update) Name() string {
	return s.fs.Name()
}

// 1. Assume the filename is the article identifier
func (s *Update) Run(cfg *Config) error {

	ctx := context.Background()

	identifier := strings.TrimSuffix(path.Base(s.filename), ".md")

	// Pull event metadata
	event, err := requestEventByIdentifier(ctx, cfg.Nsec, cfg.Relays, identifier)
	if err != nil {
		return err
	}

	// Update the event content with the new content in the identifier.md file
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}
	event.Content = string(data)

	err = publishEvent(event, cfg.Nsec, cfg.Relays)
	if err != nil {
		return err
	}

	fmt.Println("[*] event published")

	return nil
}
