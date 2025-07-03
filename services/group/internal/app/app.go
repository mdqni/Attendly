package app

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/mdqni/Attendly/services/group/internal/client"
	config2 "github.com/mdqni/Attendly/services/group/internal/config"
	"github.com/mdqni/Attendly/services/group/internal/delivery/grpc"
	"github.com/mdqni/Attendly/services/group/internal/repository/postgres"
	"github.com/mdqni/Attendly/services/group/internal/service"
	"github.com/mdqni/Attendly/shared/redisUtils"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
)

type App struct {
	server       *g.Server
	log          *slog.Logger
	address      string
	healthServer *health.Server
}

func NewApp(cfg *config2.Config, log *slog.Logger, limiter *redisUtils.Limiter, client *client.UserClient) *App {
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

	svc := service.NewGroupService(repo, cfg, client)

	server := g.NewServer(g.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
	))
	grpc.Register(
		server,
		svc,
	)
	hs := health.NewServer()
	healthpb.RegisterHealthServer(server, hs)

	return &App{server: server, log: log, address: cfg.GRPC.Address, healthServer: hs}
}

func (a *App) Run() {
	const op = "app.runGRPC"
	lis, err := net.Listen("tcp", a.address)
	if err != nil {
		panic(err)
	}
	a.healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

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
