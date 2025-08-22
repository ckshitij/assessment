package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Server struct {
	Port int `mapstructure:"port"`
}

type Backend struct {
	BaseURL  string `mapstructure:"baseurl"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Config struct {
	AppServer  Server  `mapstructure:"server"`
	NodeServer Backend `mapstructure:"backend"`
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	return &cfg
}
