package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)


type Config struct {
	// Discord
	Token       string `required:"true"`
	ReconnectDelay int `default:"5"`
	MainChannel string `required:"true"`

	// Game settings
	StreakDays int `default:"3"` // How many days in a row to count as a streak

	// SQLite
	DatabasePath string `default:"thirteenthirtyseven.db"`
}

func LoadConfig() *Config {
	cfg := &Config{}
	
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

