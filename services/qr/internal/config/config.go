package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	JwtSecret      string     `yaml:"jwtSecret" env:"JWT_SECRET" env-default:"SUPER-SECRET-CODE"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	Redis          RedisConfig `yaml:"redis"`
}

type RedisConfig struct {
	Addr    string `yaml:"addr" env-default:"redis:6379"`
	Timeout string `yaml:"timeout" env-default:"5s"`
}

type GRPCConfig struct {
	Address string `yaml:"address" env-default:"0.0.0.0:50054"`
}

func MustLoad() *Config {
	cfg := Config{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("failed to load env vars: ", err)
	}

	return &cfg
}
