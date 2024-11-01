package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const (
  configFileName = ".gatorconfig.json"
)

type Config struct {
  DBUrl            string `json:"db_url"`
  CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
  cfgFilePath, err := getConfigFilePath()
  if err != nil {
    return Config{}, err
  }

  content, err := os.ReadFile(cfgFilePath)
  if err != nil {
    return Config{}, err
  }

  config := Config{}
  if err := json.Unmarshal(content, &config); err != nil {
    return Config{}, err 
  }

  return config, nil
}

func (c *Config) SetUser(username string) error {
  if username == "" {
    return errors.New("please provide a username")
  }
  c.CurrentUserName = username
  return write(*c)
}

func getConfigFilePath() (string, error) {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    return "", err
  }

  return fmt.Sprintf("%s/%s", homeDir, configFileName), nil
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}


