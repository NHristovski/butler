package config

type Config struct {
	Server struct {
		NetworkType string `yaml:"networkType"`
		Port        string `yaml:"port"`
	} `yaml:"server"`
}
