package main

import (
	"fmt"

	"github.com/nbd-wtf/go-nostr"
)

type Note struct {
	*nostr.Event
	Path string // Full note path /tmp/zk/identifier.md
}

func (s *Note) Identifier() string {
	var res string
	for _, t := range s.Tags {
		if t.Key() == "d" {
			res = t.Value()
		}
	}
	if res == "" {
		fmt.Printf("Article with ID (%s) has no identifier\n", s.ID)
	}
	return res
}

func (s *Note) Title() string {
	var res string
	for _, t := range s.Tags {
		if t.Key() == "title" {
			res = t.Value()
		}
	}
	return res
}
