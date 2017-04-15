package main

import (
	"github.com/urfave/cli"
	"log"
)

func downCmd(c *cli.Context) error {
	log.Println("Down command call")
	return nil
}
