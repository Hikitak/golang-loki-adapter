package config

import (
	"io/ioutil"
	"os"

	"golang-loki-adapter.local/pkg/models"
	"gopkg.in/yaml.v3"
)

func LoadConfig() (*models.Config, error) {
	configFile, err := ioutil.ReadFile(os.Getenv("CONFIG_PATH"))
	if err != nil {
		return nil, err
	}

	var conf models.Config
	if err = yaml.Unmarshal(configFile, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
