package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env     string
	ConnStr string

	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := Config{
		Env:         os.Getenv("ENV"),
		ConnStr:     os.Getenv("CONSTR"),
		Address:     os.Getenv("ADDRESS"),
		Timeout:     time.Duration(time.Second * 60),
		IdleTimeout: time.Duration(time.Second * 60),
	}

	return &cfg

}
