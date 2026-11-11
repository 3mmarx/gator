package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() Config {
	configFile, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("could not locate config file: %v", err)
	}

	fileBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("could not read config file %s: %v", configFile, err)
	}

	var config Config
	if err := json.Unmarshal(fileBytes, &config); err != nil {
		log.Fatalf("could not parse config file %s: %v", configFile, err)
	}

	return config
}

func (c *Config) SetUser(username string) {
	c.CurrentUserName = username

	configFile, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("could not locate config file: %v", err)
	}

	file, err := os.Create(configFile)
	if err != nil {
		log.Fatalf("could not open config file %s: %v", configFile, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(c); err != nil {
		log.Fatalf("could not write config file %s: %v", configFile, err)
	}
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	return filepath.Join(homePath, configFileName), nil
}
