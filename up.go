package main

import (
	"log"

	"github.com/urfave/cli"
)

func upCmd(c *cli.Context) error {
	log.Println("Up command call")
	return nil
}
