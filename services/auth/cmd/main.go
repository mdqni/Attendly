package main

import (
	"context"
	"fmt"
	"github.com/mdqni/Attendly/services/auth/internal/app"
	"github.com/mdqni/Attendly/services/auth/internal/config"
	"github.com/mdqni/Attendly/services/auth/internal/kafka"
	"github.com/mdqni/Attendly/shared/redisUtils"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	rdb := redisUtils.NewRedisClient(cfg.Redis.Addr)

	limiter := redisUtils.NewLimiter(rdb)

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Error("Redis connection failed: %v", err)
	}

	fmt.Println("Redis PING OK:", pong)

	if err != nil {
		log.Error("Group client failed: %v", err)
	}

	prod, err := kafka.NewEventProducer(os.Getenv("KAFKA_BROKERS"))
	if err != nil {
		log.Error("kafka producer init error: %v", err)
	}
	defer prod.Close()

	app := app.NewApp(cfg, log, limiter, prod)
	go func() {
		app.Run()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	app.Shutdown()
	log.Info("Gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
