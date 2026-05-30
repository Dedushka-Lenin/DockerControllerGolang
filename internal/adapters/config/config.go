package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	API          ApiConfig          `json:"api"`
	DB           DatabaseConfig     `json:"database"`
	Token        TokenConfig        `json:"token"`
	Repositories RepositoriesConfig `json:"repositories"`
}

type ApiConfig struct {
	Url string `json:"url"`
}

type DatabaseConfig struct {
	DBName   string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type TokenConfig struct {
	SecretKey     string `json:"key"`
	ExpireSeconds int    `json:"expire_seconds"`
}

type RepositoriesConfig struct {
	Path string `json:"path"`
}

func LoadConfig(path string) (*Config, error) {
	conf := new(Config)

	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("LoadConfig. ReadFile. err: " + err.Error())
		return nil, fmt.Errorf("failed to read config data: %w", err)
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		log.Println("LoadConfig. Unmarshal. err: " + err.Error())
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	log.Println("LoadConfig. err: nil")
	return conf, nil
}
