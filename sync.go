package main

import (
	"log"

	"github.com/urfave/cli"
)

func syncCmd(c *cli.Context) error {
	log.Println("Sync command call")
	return nil
}
