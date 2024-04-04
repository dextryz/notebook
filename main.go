package main

import (
	"log"
	"os"
)

type Runner interface {
	Parse([]string) error
	Name() string
	Run(*Config) error
}

func main() {

	cfgPath, ok := os.LookupEnv("NOSTR")
	if !ok {
		log.Fatalln("NOSTR env var not set")
	}

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		log.Fatalln(err)
	}

	cmds := []Runner{
		NewCreate(),
		NewList(),
		NewUpdate(),
		NewPull(),
	}

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalln(err)
	}

	for _, cmd := range cmds {
		if cmd.Name() == args[0] {
			cmd.Parse(args[1:])
			err := cmd.Run(cfg)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
