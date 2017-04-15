package main

import (
	"github.com/urfave/cli"
	"log"
)

func statusCmd(c *cli.Context) error {
	log.Println("Status command call")
	return nil
}
