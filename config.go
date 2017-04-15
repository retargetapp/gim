package main

import (
	"path/filepath"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type (
	LoadConfigError error
)

type Config struct {
	Driver	string	`yaml:"driver"`
	DSN	string	`yaml:"dsn"`
	Src	string	`yaml:"src"`
}

func loadConfig () (*Config, error) {
	cfg := &Config{}
	filename, err := filepath.Abs("./.gim.yml")
	if err != nil {

		return cfg, LoadConfigError(err)
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, LoadConfigError(err)
	}

	err = yaml.Unmarshal(yamlFile, cfg)

	if err != nil {
		return cfg, LoadConfigError(err)
	}
	return cfg, nil
}
