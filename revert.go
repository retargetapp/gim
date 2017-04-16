package main

import (
	"log"

	"github.com/urfave/cli"
)

func revertCmd(c *cli.Context) error {
	log.Println("Revert command call")
	return nil
}
