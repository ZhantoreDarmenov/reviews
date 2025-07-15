package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	Database struct {
		Driver string `yaml:"driver"`
		URL    string `yaml:"url"`
	} `yaml:"database"`
}

func LoadConfig() Config {
	var cfg Config

	// Получаем путь из переменной окружения или используем по умолчанию
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}

	// Переопределяем значениями из переменных окружения, если есть
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.Database.URL = v
	}
	if v := os.Getenv("SERVER_ADDRESS"); v != "" {
		cfg.Server.Address = v
	}

	return cfg
}
