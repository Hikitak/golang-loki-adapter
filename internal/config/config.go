package config

import (
	"io/ioutil"

	"golang-loki-adapter.local/pkg/models"
	"gopkg.in/yaml.v3"
)

func LoadConfig() (*models.Config, error) {
	configFile, err := ioutil.ReadFile("internal/config/config.yaml")
	if err != nil {
		return nil, err
	}

	var conf models.Config
	if err = yaml.Unmarshal(configFile, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}