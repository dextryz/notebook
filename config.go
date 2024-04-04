package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Nsec   string   `json:"nsec"`
	Relays []string `json:"relays"`
}

func LoadConfig(path string) (*Config, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Config file: %v", err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg, nil
}
