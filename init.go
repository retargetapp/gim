package main

import (
	"github.com/urfave/cli"
	"log"
)

func initCmd(c *cli.Context) error {
	log.Println("Init command call")
	configCmd(c)
	installCmd(c)
	return nil
}
