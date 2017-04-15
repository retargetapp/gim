package main

import (
	"log"
	"fmt"
	"path/filepath"
	"os"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
	"github.com/pkg/errors"
)

func configCmd(c *cli.Context) error {
	log.Println("Config command call")
	cfg, err := core.LoadConfig()

	switch err.(type) {
	case core.InvalidConfigFileFormat:
		fmt.Println("Config file exists but it has invalid format. Setup new config:")
	case nil:
		fmt.Println("Gim already configured. Leave fields empty to use previous value:")
	}

	if err != nil {
		cfg = &core.Config{
			Driver: "mysql",
			DSN: "root:password@tcp(127.0.0.1:3306)/gim",
			Src: "./migrations",
		}
	}

	for err = errors.New(""); err != nil; {
		fmt.Printf("Database driver (%s):", cfg.Driver)
		fmt.Scanln(&cfg.Driver)

		fmt.Printf("Data source name (%s):", cfg.DSN)
		fmt.Scanln(&cfg.DSN)

		_, err = core.InitDB(cfg.Driver, cfg.DSN)
		if err != nil {
			fmt.Println("Unable connect to DB: ", err.Error())
		}
	}

	for err = errors.New(""); err != nil; {
		var path string
		var f *os.File
		fmt.Printf("Migration files sources directory (%s):", cfg.Src)
		fmt.Scanln(&cfg.Src)

		path, err = filepath.Abs(cfg.Src)
		if err != nil {
			fmt.Println("Unable to check resources directory: " + err.Error())
			continue
		}

		f, err = os.Open(path)
		defer f.Close()
		if err != nil {
			fmt.Println("Unable to check resources directory: " + err.Error())
		}
	}

	err = core.SaveConfig(cfg)

	if err != nil {
		fmt.Println("Unable to save config file: ", core.CONFIG_FILE_NAME, err.Error())
		// TODO: refactor to property exit way
		os.Exit(1)
	}
	fmt.Println("Config %s updated", core.CONFIG_FILE_NAME)
	return nil
}
