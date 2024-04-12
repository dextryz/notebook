package main

import (
	"context"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/fiatjaf/eventstore"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// 1. Notebook abstracts nostr, and therefore does not work with nostr.Event, rather with nz.Note

// Abstract working with nostr relays and local embedded database
type Notebook struct {
	cfg *Config
	db  eventstore.Store
	fs  *FileStorage
}

func NewNotebook(cfg *Config, db eventstore.Store, dir string) *Notebook {

	fs := &FileStorage{
		dir: dir,
	}

	return &Notebook{
		cfg: cfg,
		db:  db,
		fs:  fs,
	}
}

func (s *Notebook) Init() ([]*Note, error) {

	ctx := context.Background()

	// Create the directory structure on the filesystem
	err := s.fs.Init()
	if err != nil {
		return nil, err
	}

	// 1. Look inside local FS first for cached notes.

	// 2. Pull latest notes from relays.

	events, err := requestSortedEvents(context.Background(), s.cfg.Nsec, s.cfg.Relays)
	if err != nil {
		return nil, err
	}

	// 3. Store article events in local eventstore

	wdb := eventstore.RelayWrapper{Store: s.db}

	for _, e := range events {
		err := wdb.Publish(ctx, *e)
		if err != nil {
			return nil, err
		}
	}

	// 4. Store notes as markdown files on local filesystem

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

// TODO maybe move to nostr.go fix this fucking broken abstraction.
// 1. This function only connects to the local eventstore.
func (s *Notebook) Search(filter nostr.Filter) ([]*Note, error) {

	ctx := context.Background()

	var pub string
	if _, s, err := nip19.Decode(s.cfg.Nsec); err == nil {
		if pub, err = nostr.GetPublicKey(s.(string)); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

	filter.Authors = []string{pub}

	// fetch from local store if available
	wdb := eventstore.RelayWrapper{Store: s.db}

	// Try to fetch in our internal eventstore (cache) first
	events, err := wdb.QuerySync(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("no events found with search filter")
	}

	notes := []*Note{}
	for _, e := range events {
		n := &Note{Event: e}
		n.Path = fmt.Sprintf("%s/%s.md", s.fs.dir, n.Identifier())
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

// 1. Read file content from local filesystem
// 2. Pull event metadata. If event has not been pushed before, ignore
// 3. Push event with updated content and parameters
// TODO: Add a caching layer for the note object itself.
func (s *Notebook) Publish(filepath, title string, hashtags []string) error {

	content, err := s.fs.Read(filepath)
	if err != nil {
		return err
	}

	var tags nostr.Tags

	identifier := strings.TrimSuffix(path.Base(filepath), ".md")

	tags = append(tags, nostr.Tag{"d", identifier})
	tags = append(tags, nostr.Tag{"title", title})

	for _, v := range hashtags {
		tags = append(tags, nostr.Tag{"t", v})
	}

	e := &nostr.Event{
		Kind:      nostr.KindArticle,
		Content:   string(content),
		CreatedAt: nostr.Now(),
		Tags:      tags,
	}

	// 1. Publish event to nostr relays for distributed/global storage

	// TODO maybe remove this and use Push expliciteyly
	err = publishEvent(e, s.cfg.Nsec, s.cfg.Relays)
	if err != nil {
		return err
	}

	fmt.Printf("[*] event published with ID: %s\n", e.ID)

	return nil
}

// 1. Store note on local file system for caching/local storage
func (s *Notebook) Add() (*Note, error) {

	var tags nostr.Tags

	identifier := time.Now().Format("200601021504")
	tags = append(tags, nostr.Tag{"d", identifier})

	e := &nostr.Event{
		Kind:      nostr.KindArticle,
		Content:   "",
		CreatedAt: nostr.Now(),
		Tags:      tags,
	}

	n := &Note{Event: e}

	err := s.fs.Store(n)
	if err != nil {
		return nil, err
	}

	return n, nil
}
