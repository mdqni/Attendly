package main

import (
	"context"
	"fmt"
	app2 "github.com/mdqni/Attendly/services/group/internal/app"
	"github.com/mdqni/Attendly/services/group/internal/client"
	"github.com/mdqni/Attendly/services/group/internal/config"
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

	app2.RunMigrations(cfg.ConnString, "/app/internal/migrations")

	log := setupLogger(cfg.Env)

	rdb := redisUtils.NewRedisClient(cfg.Redis.Addr)
	limiter := redisUtils.NewLimiter(rdb)
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Error("Redis connection failed: %v", err)
	}
	fmt.Println("Redis PING OK:", pong)

	userClient, err := client.NewUserClient(cfg.UserServiceAddr)
	if err != nil {
		log.Error("failed to connect to user service: %v", err)
	}
	defer func(userClient *client.UserClient) {
		err := userClient.Close()
		if err != nil {
			log.Error("failed to close user client: %v", err)
		}
	}(userClient)

	app := app2.NewApp(cfg, log, limiter, userClient)
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
