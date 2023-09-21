package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"

	"log"
	"os"
)

var config Config

func init() {
	err := godotenv.Load("./config/config.yaml")
	e := os.Getenv("env")
	if e != "production" {
		if err != nil {
			log.Fatal("Error on load configuration file.")
		}
	}

	if err := env.Parse(&config); err != nil {
		log.Fatal("Error on parsing configuration file.", err)
	}

	log.Printf(`
		env: %s`,
		config.Environment,
	)
}

// GetConfig : getter
func GetConfig() *Config {
	return &config
}

type Config struct {
	Environment string `json:"env" env:"env"`
	Port        string `json:"port" env:"port"`
	DB          string `json:"db" env:"db"`
	DBName      string `json:"db_name" env:"db_name"`
}
