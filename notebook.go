package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

type FileStorage struct {
	dir string
}

// 1. Notebook abstracts nostr, and therefore does not work with nostr.Event, rather with nz.Note

// Abstract working with nostr relays and local embedded database
type Notebook struct {
	cfg *Config
	fs  *FileStorage
}

func NewNotebook(cfg *Config, dir string) *Notebook {

	// TODO: Handle the case for when dir already exists
	err := os.Mkdir(dir, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	fs := &FileStorage{
		dir: dir,
	}

	return &Notebook{
		cfg: cfg,
		fs:  fs,
	}
}

func (s *Notebook) Init() ([]*Note, error) {

	// 1. Look inside local FS first for cached notes.

	// 2. Pull latest notes from relays.

	events, err := requestSortedEvents(context.Background(), s.cfg.Nsec, s.cfg.Relays)
	if err != nil {
		return nil, err
	}

	notes := []*Note{}
	for _, e := range events {

		n := &Note{Event: e}

		// Store on local filesystem
		err = s.fs.Store(n)
		if err != nil {
			return nil, err
		}

		notes = append(notes, n)
	}

	return notes, nil
}

func (s *Notebook) FindNotes() ([]*Note, error) {

	// 1. Look inside local FS first for cached notes.

	// 2. Pull latest notes from relays.

	events, err := requestSortedEvents(context.Background(), s.cfg.Nsec, s.cfg.Relays)
	if err != nil {
		return nil, err
	}

	notes := []*Note{}
	for _, e := range events {
		n := &Note{Event: e}
		notes = append(notes, n)
	}

	return notes, nil
}

func (s *Notebook) Add(title, content string, hashtags []string) (*Note, error) {

	var tags nostr.Tags

	identifier := time.Now().Format("200601021504")
	tags = append(tags, nostr.Tag{"d", identifier})
	tags = append(tags, nostr.Tag{"title", title})

	for _, v := range hashtags {
		tags = append(tags, nostr.Tag{"t", v})
	}

	e := &nostr.Event{
		Kind:      nostr.KindArticle,
		Content:   content,
		CreatedAt: nostr.Now(),
		Tags:      tags,
	}

	// 1. Publish event to nostr relays for distributed/global storage

	err := publishEvent(e, s.cfg.Nsec, s.cfg.Relays)
	if err != nil {
		return nil, err
	}

	// 2. Store note on local file system for caching/local storage

	n := &Note{Event: e}

	err = s.fs.Store(n)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// 1. Store event in local working directory
// 2. Store event in SQLite for fast querying
func (s *FileStorage) Store(n *Note) error {

	filename := fmt.Sprintf("%s/%s.md", s.dir, n.Identifier())

	err := os.WriteFile(filename, []byte(n.Content), 0660)
	if err != nil {
		return err
	}

	return nil
}
