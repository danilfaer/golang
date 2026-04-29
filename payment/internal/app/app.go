package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/danilfaer/golang/payment/internal/config"
	"github.com/danilfaer/golang/platform/pkg/closer"
	"github.com/danilfaer/golang/platform/pkg/logger"
	sharedInterceptors "github.com/danilfaer/golang/shared/pkg/interceptors"
	paymentV1 "github.com/danilfaer/golang/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// App — composition root payment-сервиса (gRPC).
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

// Run блокируется на Serve до остановки (graceful shutdown через closer).
func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	// Порядок важен для обратного закрытия в closer (последний зарегистрированный — первым).
	inits := []func(ctx context.Context) error{
		a.initDiContainer,
		a.initLogger,
		a.initCloser,
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

func (a *App) initDiContainer(ctx context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

// initLogger поднимает zap по настройкам из .env (уровень, JSON).
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

	paymentV1.RegisterPaymentServiceServer(s, a.diContainer.PaymentGRPCAPI(ctx))
	reflection.Register(s)

	a.grpcServer = s

	// GracefulStop разблокирует Serve после сигнала (через closer).
	closer.AddNamed("gRPC Server", func(ctx context.Context) error {
		// Блокируется до завершения активных RPC или отмены контекста shutdown.
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
