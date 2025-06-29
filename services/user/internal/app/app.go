package app

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/mdqni/Attendly/services/user/internal/config"
	"github.com/mdqni/Attendly/services/user/internal/delivery/grpc"
	"github.com/mdqni/Attendly/services/user/internal/interceptor"
	"github.com/mdqni/Attendly/services/user/internal/repository/postgres"
	"github.com/mdqni/Attendly/services/user/internal/service"
	"github.com/mdqni/Attendly/shared/redisLimiter"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
)

type App struct {
	server  *g.Server
	log     *slog.Logger
	address string
}

func NewApp(cfg *config.Config, log *slog.Logger, address string, limiter *redisLimiter.Limiter) *App {
	const op = "app.NewApp"

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	repo, err := postgres.New(cfg.ConnString)
	if err != nil {
		log.Error("failed to init postgres", slog.String("op", op), slog.Any("err", err))
		panic(err)
	}

	svc := service.NewUserService(repo, limiter, cfg)

	server := g.NewServer(g.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		interceptor.RBACInterceptor(cfg.JwtSecret),
	))
	grpc.Register(
		server,
		svc,
	)

	return &App{server: server, log: log, address: address}
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
