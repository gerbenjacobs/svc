package internal

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	Svc struct {
		Name        string
		Version     string
		Env         string
		Address     string
		SecretToken string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
	}
}

func NewConfig() Configuration {
	// load configuration
	var c Configuration
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}

	return c
}
