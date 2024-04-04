package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// 1.
// > nz new 202402051756.md --title "Fake Knowledge" --tag nostr --tag bitcoin

func main() {

	ctx := context.Background()

	cfgPath, ok := os.LookupEnv("NOSTR")
	if !ok {
		log.Fatalln("NOSTR env var not set")
	}

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		log.Fatalln(err)
	}

	switch os.Args[1] {
	case "new":
		fs := flag.NewFlagSet("new", flag.ExitOnError)
		e := buildEvent(fs, os.Args[2:])
		publishEvent(cfg, e)
	case "list":
		events, _ := Request(ctx, cfg)
		ViewTitle(events)
	case "pull":
		events, _ := Request(ctx, cfg)
		Store(events, os.Args[2])
	case "update":

		filename := os.Args[2]

		identifier := strings.TrimSuffix(path.Base(filename), ".md")

		// Pull event metadata
		event, err := RequestByIdentifier(ctx, cfg, identifier)
		if err != nil {
			log.Fatal(err)
		}

		// Update the event content with the new content in the identifier.md file
		data, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		event.Content = string(data)

		publishEvent(cfg, event)

		fmt.Println("[*] event published")

	default:
		fmt.Println("ERROR")
	}
}

type Article struct {
	identifier string
	title      string
	content    string
	tags       []string
}

func (s Article) Print() {
	fmt.Printf("%s\n", s.identifier)
	fmt.Printf("- %s\n", s.title)
	fmt.Println("")
}

func (s Article) Publish() {
	fmt.Println("publishing to nostr")
	fmt.Println(s)
}

// Build event data structure from commandline arguments
// 1. Read content from input file
// 2. Create list of tags
// 3. Add title
// 4. Create unique identifier
func buildEvent(fs *flag.FlagSet, args []string) *nostr.Event {

	filename := args[0]
	if path.Ext(filename) != ".md" {
		log.Fatalln("file type has to be .md")
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
    content := string(data)

	var tags nostr.Tags
	fs.Func("tag", "tag article", func(s string) error {
		tags = append(tags, nostr.Tag{"t", s})
		return nil
	})

    var title string
	fs.StringVar(&title, "title", "", "the filename to be committed")
	tags = append(tags, nostr.Tag{"title", title})

    identifier := time.Now().Format("200601021504")
	tags = append(tags, nostr.Tag{"d", identifier})

	fs.Parse(args)

	return &nostr.Event{
		Kind:      nostr.KindArticle,
		Content:   content,
		CreatedAt: nostr.Now(),
		Tags:      tags,
	}
}

func RequestByIdentifier(ctx context.Context, cfg *Config, identifier string) (*nostr.Event, error) {

	var pub string
	if _, s, err := nip19.Decode(cfg.Nsec); err == nil {
		if pub, err = nostr.GetPublicKey(s.(string)); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	filter := nostr.Filter{
		Kinds:   []int{nostr.KindArticle},
		Authors: []string{pub},
		Tags: nostr.TagMap{
			"d": []string{identifier},
		},
		Limit: 1,
	}

	events := cfg.queryRelays(ctx, filter)

	if len(events) == 0 {
		return nil, fmt.Errorf("not article found with identifier: %s", identifier)
	}

	// Pop the lastest parameterized replaceable event
	return events[len(events)-1], nil
}

func Request(ctx context.Context, cfg *Config) ([]*nostr.Event, error) {

	var pub string
	if _, s, err := nip19.Decode(cfg.Nsec); err == nil {
		if pub, err = nostr.GetPublicKey(s.(string)); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	filter := nostr.Filter{
		Kinds:   []int{nostr.KindArticle},
		Authors: []string{pub},
		Limit:   500,
	}

	events := cfg.queryRelays(ctx, filter)

	slices.SortFunc(events, func(a, b *nostr.Event) int { return int(b.CreatedAt - a.CreatedAt) })

	return events, nil
}

func ViewTitle(events []*nostr.Event) {

	for _, e := range events {

		a := Article{}

		for _, t := range e.Tags {
			if t.Key() == "title" {
				a.title = t.Value()
			}
			if t.Key() == "d" {
				a.identifier = t.Value()
			}
		}

		if a.identifier == "" {
			fmt.Printf("Article with title (%s) has no identifier\n", a.title)
		} else {
			a.Print()
		}
	}
}

func Store(events []*nostr.Event, path string) error {

	err := os.Mkdir(path, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	for _, e := range events {

		a := Article{}

		for _, t := range e.Tags {
			if t.Key() == "d" {
				a.identifier = t.Value()
			}
		}

		if a.identifier == "" {
			fmt.Printf("Article with title (%s) has no identifier\n", a.title)
		} else {
			filename := fmt.Sprintf("%s/%s.md", path, a.identifier)
			err = os.WriteFile(filename, []byte(e.Content), 0660)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 1. Add author pubkey
// 2. Add created at timestamp
// 3. Sign event
// 4. Publish to set of relays
func publishEvent(cfg *Config, e *nostr.Event) error {

	ctx := context.Background()

	var sk string
	var pub string
	if _, s, err := nip19.Decode(cfg.Nsec); err == nil {
		sk = s.(string)
		if pub, err = nostr.GetPublicKey(s.(string)); err != nil {
            log.Fatal(err)
		}
	} else {
        log.Fatal(err)
	}

    e.PubKey = pub
    e.CreatedAt = nostr.Now()

	err := e.Sign(sk)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, r := range cfg.Relays {
		wg.Add(1)

		go func(relayUrl string) {
			defer wg.Done()

			relay, err := nostr.RelayConnect(ctx, relayUrl)
			if err != nil {
				log.Println(err)
				return
			}
			defer relay.Close()

			err = relay.Publish(ctx, *e)
			if err != nil {
				log.Println(err)
				return
			}
		}(r)
	}
	wg.Wait()

	return nil
}
