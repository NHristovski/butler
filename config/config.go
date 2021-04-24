package config

type Config struct {
	Server struct {
		NetworkType string `yaml:"networkType"`
		Port        string `yaml:"port"`
		RootDir     string `yaml:"rootDir"`
	} `yaml:"server"`

	Logging struct {
		Debug struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"debug"`
	} `yaml:"logging"`
}
