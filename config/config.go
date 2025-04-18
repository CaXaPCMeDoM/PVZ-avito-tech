package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Jwt        JWT
		HTTP       HTTP
		Log        Log
		Pg         PG
		Security   Security
		Prometheus Prometheus
	}

	JWT struct {
		SecretKey string `env:"JWT_SECRET_KEY" env-required:"true"`
	}

	HTTP struct {
		Port            string        `env:"HTTP_PORT" env-required:"true"`
		ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"30s"`
		WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"30s"`
		IdleTimeout     time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
		MaxHeaderBytes  int           `env:"HTTP_MAX_HEADER_BYTES" env-default:"1048576"`
		Mode            string        `env:"GIN_MODE" env-required:"true"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"5s"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-required:"true"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX" env-required:"true"`
		URL     string `env:"PG_URL" env-required:"true"`
	}

	Security struct {
		PasswordCost int `env:"SECURITY_PASSWORD_COST" env-default:"10"`
	}

	Prometheus struct {
		Enabled bool   `env:"METRICS_ENABLED" env-required:"true"`
		Port    string `env:"METRICS_PORT" env-required:"true"`
	}
)

func MustLoad() *Config {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("error loading .env file: %v", err)
		}
	}

	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	if cfg.Security.PasswordCost > 12 {
		log.Fatal("SECURITY_PASSWORD_COST is too high. It should be <13")
	}
	if cfg.HTTP.ReadTimeout < 0 {
		log.Fatal("HTTP_READ_TIMEOUT cannot be negative")
	}
	if cfg.HTTP.WriteTimeout < 0 {
		log.Fatal("HTTP_WRITE_TIMEOUT cannot be negative")
	}

	return &cfg
}
