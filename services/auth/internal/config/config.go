package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type Config struct {
	Env             string     `yaml:"env" env-default:"local"`
	ConnString      string     `yaml:"connString" env:"CONN_STRING" env-default:"postgresql://postgres:postgres@localhost:5432/attendly?sslmode=disable"`
	JwtSecret       string     `yaml:"jwtSecret" env:"JWT_SECRET" env-default:"SUPER-SECRET-CODE"`
	GRPC            GRPCConfig `yaml:"grpc"`
	MigrationsPath  string
	JwtTokenTTL     time.Duration `yaml:"jwt_token_ttl" env-default:"1h"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-default:"168h"`
	Redis           RedisConfig   `yaml:"redis"`
}

type RedisConfig struct {
	Addr string `yaml:"addr" env-default:"redis:6379"`
}

type GRPCConfig struct {
	Address string        `yaml:"address" env-default:"0.0.0.0:50050"`
	Timeout time.Duration `yaml:"timeout" env-default:"1h"`
}

func MustLoad() *Config {
	cfg := Config{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("failed to load env vars: ", err)
	}

	return &cfg
}
