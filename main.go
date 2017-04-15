package main

import (
	"log"
	"fmt"
)

func main() {
	log.Println("Run")
}

func init() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Unable to load config from .gim.yml: %v\n", err)
	}
	log.Printf("Config: %#v\n", cfg)

	err = initDB(cfg.Driver, cfg.DSN)

	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
	}

	m, _ :=loadSources(cfg.Src)
	log.Printf("%#v\n", m[0])
	// Check db
	// – Create DB migrations table

	// Load migrations
	//
}
