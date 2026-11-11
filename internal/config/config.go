package config

import (
	"encoding/json"
	"errors"
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
		log.Fatalf("Failed to get config file: %s", err)
	}

	fileBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read the config file: %s", err)
	}

	var config Config
	if err := json.Unmarshal(fileBytes, &config); err != nil {
		log.Fatalf("Failed to unmarshal json: %s", err)
	}

	return config
}

func (c *Config) SetUser(username string) {
	c.CurrentUserName = username

	configFile, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("Failed to get config file: %s", err)
	}

	file, err := os.Create(configFile)
	if err != nil {
		log.Fatalf("Failed to open the config file: %s", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(c); err != nil {
		log.Fatalf("failed to encode config JSON: %s", err)
	}
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("can't find home path")
	}

	return filepath.Join(homePath, configFileName), nil
}
