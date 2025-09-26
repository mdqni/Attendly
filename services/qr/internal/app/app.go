package app

import (
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/mdqni/Attendly/services/qr/internal/config"
	"github.com/mdqni/Attendly/services/qr/internal/delivery/grpc"
	redis2 "github.com/mdqni/Attendly/services/qr/internal/repository/redis"
	"github.com/mdqni/Attendly/services/qr/internal/service"
	"github.com/redis/go-redis/v9"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	server  *g.Server
	log     *slog.Logger
	address string
}

func NewApp(cfg *config.Config, log *slog.Logger, rdb *redis.Client) *App {
	const op = "app.NewApp"

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	cache := redis2.NewRedisCache(rdb)

	svc := service.NewQrService(cache)

	server := g.NewServer(g.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
	))
	grpc.Register(
		server,
		svc,
	)

	return &App{server: server, log: log, address: cfg.GRPC.Address}
}

func (a *App) Run() {
	const op = "app.runGRPC"
	lis, err := net.Listen("tcp", a.address)
	if err != nil {
		panic(err)
	}

	a.log.Info("grpc server started", "op:", op, slog.String("addr", lis.Addr().String()))
	if err := a.server.Serve(lis); err != nil {
		panic(err)
	}
}

func (a *App) Shutdown() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.String("address", a.address))

	a.server.GracefulStop()
}
