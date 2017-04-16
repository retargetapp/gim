package main

import (
	"log"

	"github.com/urfave/cli"
)

func statusCmd(c *cli.Context) error {
	log.Println("Status command call")
	return nil
}
