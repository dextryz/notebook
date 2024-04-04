package main

import (
	"context"
	"fmt"
	"log"
	"slices"
	"sync"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// 1. Add author pubkey
// 2. Add created at timestamp
// 3. Sign event
// 4. Publish to set of relays
func publishEvent(e *nostr.Event, nsec string, relays []string) error {

	ctx := context.Background()

	var sk string
	var pub string
	if _, s, err := nip19.Decode(nsec); err == nil {
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
	for _, r := range relays {
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

// Last events are last. This ensures that doublicates are overwritten when stored.
func requestSortedEvents(ctx context.Context, nsec string, relays []string) ([]*nostr.Event, error) {

	var pub string
	if _, s, err := nip19.Decode(nsec); err == nil {
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

	events := queryRelays(ctx, filter, relays)

	slices.SortFunc(events, func(a, b *nostr.Event) int { return int(a.CreatedAt - b.CreatedAt) })

	return events, nil
}

func requestEventByIdentifier(ctx context.Context, nsec string, relays []string, identifier string) (*nostr.Event, error) {

	var pub string
	if _, s, err := nip19.Decode(nsec); err == nil {
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

	events := queryRelays(ctx, filter, relays)

	if len(events) == 0 {
		return nil, fmt.Errorf("not article found with identifier: %s", identifier)
	}

	// Pop the lastest parameterized replaceable event
	return events[len(events)-1], nil
}

func queryRelays(ctx context.Context, filter nostr.Filter, relays []string) (ev []*nostr.Event) {

	var m sync.Map
	var wg sync.WaitGroup
	for _, url := range relays {

		wg.Add(1)
		go func(wg *sync.WaitGroup, url string) {
			defer wg.Done()

			r, err := nostr.RelayConnect(ctx, url)
			if err != nil {
				panic(err)
			}

			events, err := r.QuerySync(ctx, filter)
			if err != nil {
				// TODO
				return
			}

			for _, e := range events {
				m.Store(e.ID, e)
			}

		}(&wg, url)
	}
	wg.Wait()

	m.Range(func(_, v any) bool {
		ev = append(ev, v.(*nostr.Event))
		return true
	})

	return ev
}
