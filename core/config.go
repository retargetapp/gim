package core

import (
	"path/filepath"
	"io/ioutil"
	"os"
	"gopkg.in/yaml.v2"
)

const CONFIG_FILE_NAME = ".gim.yml"

type (
	LoadConfigFileError error
	InvalidConfigFileFormat error
	SaveConfigFileError error
)

type Config struct {
	Driver	string	`yaml:"driver"`
	DSN	string	`yaml:"dsn"`
	Src	string	`yaml:"src"`
}

func LoadConfig () (*Config, error) {
	cfg := &Config{}
	filename, err := filepath.Abs(CONFIG_FILE_NAME)
	if err != nil {
		return cfg, LoadConfigFileError(err)
	}
	ymlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, LoadConfigFileError(err)
	}

	err = yaml.Unmarshal(ymlFile, cfg)

	if err != nil {
		return cfg, InvalidConfigFileFormat(err)
	}
	return cfg, nil
}

func SaveConfig(cfg *Config) error {
	yml, _ := yaml.Marshal(cfg)
	err := ioutil.WriteFile(CONFIG_FILE_NAME, yml, os.FileMode(int(0640)))
	if err != nil {
		return SaveConfigFileError(err)
	}
	return nil
}
