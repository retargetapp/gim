package main

import (
	"github.com/urfave/cli"
	"log"
)

func upCmd(c *cli.Context) error {
	log.Println("Up command call")
	return nil
}
