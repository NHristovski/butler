package config_parser

import (
	"butler/config"
	"gopkg.in/yaml.v2"
	"os"
)

func InitConfig(cfg *config.Config, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}

	return nil
}
