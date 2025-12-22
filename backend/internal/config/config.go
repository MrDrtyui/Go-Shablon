package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env"`
	Http     `yaml:"http"`
	Database `yaml:"database"`
	Jwt      `yaml:"jwt"`
}

type Http struct {
	Port string `yaml:"port"`
}

type Database struct {
	Url string `yaml:"url"`
}

type Jwt struct {
	Secret          string        `yaml:"secret"`
	AccessTtlHours  time.Duration `yaml:"accessTtlHours"`
	RefreshTtlHours time.Duration `yaml:"refreshTtlHours"`
}

func MustLoad() *Config {
	var cfg Config

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH IS NIL")
	}

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatal("Failed init configs", err.Error())
	}

	return &cfg
}
