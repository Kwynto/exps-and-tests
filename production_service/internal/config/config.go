package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEV" env-default:"false"`
	Listen        struct {
		Type   string `env:"LISTEN_TYPE" env-default:"port"`
		BindIP string `env:"BIND_IP" env-default:"0.0.0.0"`
		Port   string `env:"PORT" env-default:"10000"`
	}
	AppConfig struct {
		LogLevel  string
		AdminUser struct {
			Email    string `env:"ADMIN_EMAIL" env-required:"true"`
			Password string `env:"ADMIN_PWD" env-required:"true"`
		}
	}
}

var instance *Config

var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Print("gether config")

		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			helptext := "Monolit Notes System"
			help, _ := cleanenv.GetDescription(instance, &helptext)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
