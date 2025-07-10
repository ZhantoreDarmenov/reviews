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

	data, err := os.ReadFile("C:\\Users\\User\\Desktop\\reviews\\config\\config.yaml")
	//path := os.Getenv("CONFIG_PATH")
	//if path == "" {
	//	path = "config/config.yaml"
	//}
	//
	//data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	//
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		//if err = yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}
	//if v := os.Getenv("DATABASE_URL"); v != "" {
	//	cfg.Database.URL = v
	//}
	//
	//if v := os.Getenv("SERVER_ADDRESS"); v != "" {
	//	cfg.Server.Address = v
	//}
	return cfg
}
