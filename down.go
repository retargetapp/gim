package main

import (
	"log"

	"github.com/urfave/cli"
)

func downCmd(c *cli.Context) error {
	log.Println("Down command call")
	return nil
}
