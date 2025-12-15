package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type HTTPConfig struct {
	ListenHost string `yaml:"listen_host"`
	ListenPort int    `yaml:"listen_port"`
	Timeout    int    `yaml:"timeout"`
}

type Config struct {
	HTTPConfig HTTPConfig `yaml:"http_config"`
}

func Load() (*Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		log.Println("Error getting executable path:", err)
		return nil, err
	}

	exeDir := filepath.Dir(exePath)
	path := filepath.Join(exeDir, "config.yaml")

	// Try to load config.yaml relative to the binary
	res, err := os.ReadFile(path)
	if err != nil {
		log.Println("Error reading config from binary dir:", err)
		log.Println("Trying to read the config from current directory (pwd)")

		// Fallback to read config from pwd (for dev with go run ...)
		res, err = os.ReadFile("config.yaml")
		if err != nil {
			log.Println("Error reading config from pwd:", err)
			return nil, err
		}
	}

	var cfg *Config
	if err := yaml.Unmarshal(res, &cfg); err != nil {
		log.Println("Error reading YAML:", err)
		return nil, err
	}

	return cfg, nil
}
