package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/danilfaer/golang/inventory/internal/config"
	"github.com/danilfaer/golang/platform/pkg/closer"
	"github.com/danilfaer/golang/platform/pkg/logger"
	sharedInterceptors "github.com/danilfaer/golang/shared/pkg/interceptors"
	inventoryV1 "github.com/danilfaer/golang/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// App — composition root inventory-сервиса (gRPC + MongoDB).
type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

// New собирает зависимости и возвращает готовое приложение.
func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize dependencies: %w", err)
	}

	return a, nil
}

// Run блокируется на gRPC Serve до graceful shutdown.
func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initDIContainer,
		a.initLogger,
		a.initCloser,
		a.initMongo,
		a.initListener,
		a.initGRPCServer,
	}

	for _, init := range inits {
		if err := init(ctx); err != nil {
			return fmt.Errorf("failed to initialize %s: %w", init, err)
		}
	}

	return nil
}

func (a *App) initDIContainer(ctx context.Context) error {
	a.diContainer = newDIContainer()
	return nil
}

// initLogger поднимает zap через platform/logger по переменным из конфига.
func (a *App) initLogger(ctx context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJSON(),
	)
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

// initMongo создаёт соединение, проверяет ping и регистрирует Disconnect в closer.
func (a *App) initMongo(ctx context.Context) error {
	if err := a.diContainer.InitMongo(ctx); err != nil {
		return fmt.Errorf("mongo connect: %w", err)
	}

	c := a.diContainer.MongoClient()
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := c.Ping(pingCtx, nil); err != nil {
		return fmt.Errorf("mongo ping: %w", err)
	}

	closer.AddNamed("MongoDB", func(shutdownCtx context.Context) error {
		return c.Disconnect(shutdownCtx)
	})

	return nil
}

func (a *App) initListener(ctx context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().GRPC.Address())
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	closer.AddNamed("TCP Listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(sharedInterceptors.UnaryErrorInterceptor()),
	)

	inventoryV1.RegisterInventoryServiceServer(s, a.diContainer.InventoryGRPCAPI(ctx))
	reflection.Register(s)

	a.grpcServer = s

	closer.AddNamed("gRPC Server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("gRPC server listening on %s", config.AppConfig().GRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}

// GRPCServer — доступ к серверу (тесты/расширения).
func (a *App) GRPCServer() *grpc.Server {
	return a.grpcServer
}
