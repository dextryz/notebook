package main

import (
	"flag"
	"fmt"
	"path"
)

// > nz push /tmp/zk/202404041455.md

type Push struct {
	fs *flag.FlagSet

	filename string
	title    string
	hashtags []string
}

func NewPush() *Push {

	s := Push{
		fs: flag.NewFlagSet("push", flag.ExitOnError),
	}

	s.fs.StringVar(&s.filename, "content", "", "the filename to be committed")
	s.fs.StringVar(&s.title, "title", "", "the filename to be committed")

	s.fs.Func("tag", "tag article", func(v string) error {
		s.hashtags = append(s.hashtags, v)
		return nil
	})

	return &s
}

func (s *Push) Parse(args []string) error {
	return s.fs.Parse(args)
}

func (s *Push) Name() string {
	return s.fs.Name()
}

// 1. Assume the filename is the article identifier
func (s *Push) Run(container *Container) error {

	// Get the current notebook instance
	nb, err := container.CurrentNotebook()
	if err != nil {
		return err
	}

	// TODO check different filetyopes
	if path.Ext(s.filename) != ".md" {
		return fmt.Errorf("file type has to be .md")
	}

	err = nb.Publish(s.filename, s.title, s.hashtags)
	if err != nil {
		return err
	}

	return nil
}
