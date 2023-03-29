package config

import (
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/types"
	"go.deanishe.net/env"
	"strings"
)

const (
	ProductionEnv types.ENVIRONMENT = "production"
	SandboxEnv    types.ENVIRONMENT = "sandbox"
)

var Instance *Config

type Config struct {
	AppName           string `env:"APP_NAME"`
	Namespace         string `env:"NAMESPACE"`
	Host              string `env:"HOST"`
	Env               string `env:"ENVIRONMENT"`
	SecretKey         string `env:"SECRET_KEY"`
	TOKENLIFESPAN     uint   `env:"TOKEN_LIFE_SPAN"`
	DbName            string `env:"DB_NAME"`
	DbUser            string `env:"DB_USER"`
	DbPass            string `env:"DB_PASS"`
	DbHost            string `env:"DB_HOST"`
	DbPort            uint   `env:"DB_PORT"`
	DbURL             string `env:"DB_URL"`
	RedisURL          string `env:"REDIS_URL"`
	RedisPassword     string `env:"REDIS_PASSWORD"`
	RedisUsername     string `env:"REDIS_USERNAME"`
	SendGridAPIKey    string `env:"SENDGRID_API_KEY"`
	SendGridFromEmail string `env:"SEND_GRID_FROM_EMAIL"`
}

func Load() {
	c := &Config{}
	if err := env.Bind(c); err != nil {
		panic(err.Error())
	}
	Instance = c
	return
}

func (c *Config) GetEnv() types.ENVIRONMENT {
	return types.ENVIRONMENT(strings.ToLower(Instance.Env))
}
