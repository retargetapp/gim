package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/pkg/errors"
)

const (
	CONFIG_FILE_NAME = ".gim.yml"

	ERROR_LOAD_CONFIG_FILE = "load_config_file_error"
	ERROR_INVALID_CONFIG_FILE_FORMAT = "invalid_config_file_format"
	ERROR_SAVE_CONFIG_FILE = "save_confgi_file_error"
)

type Config struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
	Src    string `yaml:"src"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	filename, err := filepath.Abs(CONFIG_FILE_NAME)
	if err != nil {
		return cfg, errors.New(ERROR_LOAD_CONFIG_FILE)
	}
	ymlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, errors.New(ERROR_LOAD_CONFIG_FILE)
	}

	err = yaml.Unmarshal(ymlFile, cfg)

	if err != nil {
		return cfg, errors.New(ERROR_INVALID_CONFIG_FILE_FORMAT)
	}
	return cfg, nil
}

func SaveConfig(cfg *Config) error {
	yml, _ := yaml.Marshal(cfg)
	err := ioutil.WriteFile(CONFIG_FILE_NAME, yml, os.FileMode(int(0640)))
	if err != nil {
		return errors.New(ERROR_SAVE_CONFIG_FILE)
	}
	return nil
}
