package main

import (
	"log"
	"os"

	eventstore_badger "github.com/fiatjaf/eventstore/badger"
)

// Business Invarients
//
// 1. A notebook has to be created before notes can be created or pulled.
// 2. Pull into a directory creates a notebook.
// 3. Creating a new note also creates a notebook.

type Runner interface {
	Parse([]string) error
	Name() string
	Run(*Container) error
}

func main() {

	cfg, err := LoadConfig(os.Getenv("NOSTR"))
	if err != nil {
		log.Fatalln(err)
	}

	db := &eventstore_badger.BadgerBackend{
		Path: "eventstore.db",
	}
	err = db.Init()
	if err != nil {
		panic(err)
	}

	// Since we want multiple commands to set a notebook we need
	// general purpose container to check the current notebook state, the dtate of my editor and shell
	container := NewContainer(cfg, db)

	cmds := []Runner{
		NewInit(),
		NewCreate(),
		NewPush(),
		NewList(),
		NewSearch(),
	}

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalln(err)
	}

	for _, cmd := range cmds {
		if cmd.Name() == args[0] {
			cmd.Parse(args[1:])
			err := cmd.Run(container)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
