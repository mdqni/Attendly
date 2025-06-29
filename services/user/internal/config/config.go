package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	ConnString     string     `env:"CONN_STRING" env-default:"postgresql://postgres:postgres@localhost:5432/attendly?sslmode=disable"`
	JwtSecret      string     `env:"JWT_SECRET" env-default:"SUPER-SECRET-CODE"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
	Redis          RedisConfig   `yaml:"redis"`
}

type RedisConfig struct {
	Addr string `yaml:"addr" env-default:"localhost:6379"`
}

type GRPCConfig struct {
	Address string        `yaml:"address" env-default:"localhost:50051"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		configPath = "./config/config.yaml"
	}

	cfg := Config{}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("config file not found, using defaults")
		return &cfg
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config parse error: " + err.Error())
	}

	return &cfg
}

var once sync.Once

func fetchConfigPath() string {
	var res string

	once.Do(func() {
		flag.StringVar(&res, "config", "", "path to config file")
		flag.Parse()
	})

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
