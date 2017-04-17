package main

import (
	"log"

	"github.com/urfave/cli"
)

func createCmd(c *cli.Context) error {
	log.Println("Create command call")
	return nil
}
