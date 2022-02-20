package common

import (
	"os"

	"gopkg.in/yaml.v3"
)

var AppConfig Config

type Config struct {
	Kucoin struct {
		ApiKey          string `yaml:"apiKey"`
		ApiSecret       string `yaml:"apiSecret"`
		ApiPassphrase   string `yaml:"apiPassphrase"`
		SliceSize       int    `yaml:"sliceSize"`
		LendingCurrency string `yaml:"lendingCurrency"`
		Term            int    `yaml:"term"`
	}

	Poloniex struct {
		ApiKey          string `yaml:"apiKey"`
		ApiSecret       string `yaml:"apiSecret"`
		SliceSize       int    `yaml:"sliceSize"`
		LendingCurrency string `yaml:"lendingCurrency"`
	}
}

func init() {
	currentPath, _ := os.Getwd()
	file, _ := os.ReadFile(currentPath + "/config.yaml")
	yaml.Unmarshal(file, &AppConfig)
}
