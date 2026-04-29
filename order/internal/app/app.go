package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/danilfaer/golang/order/internal/config"
	"github.com/danilfaer/golang/platform/pkg/closer"
	"github.com/danilfaer/golang/platform/pkg/logger"
	order_v1 "github.com/danilfaer/golang/shared/pkg/api/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


type App struct {
	diContainer *diContainer
	apiServer *http.Server
	listener net.Listener
}


func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize dependencies: %w", err)
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initDiContainer,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initMigrations,
		a.initHTTPServer,
	}

	for _, init := range inits {
		if err := init(ctx); err != nil {
			return fmt.Errorf("failed to initialize %s: %w", init, err)
		}
	}

	return  nil
}


func (a *App) initDiContainer(ctx context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

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
	listener, err := net.Listen("tcp", config.AppConfig().HTTP.Address())
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

func (a *App) initHTTPServer(ctx context.Context) error {
	// Создаем OpenAPI сервер
	api := a.diContainer.OrderV1API(ctx)
	s, err := order_v1.NewServer(api)
	if err != nil {
		return fmt.Errorf("failed to create OpenAPI server: %w", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", s)

	a.apiServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	closer.AddNamed("HTTP Server", func(ctx context.Context) error {
		return a.apiServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) initMigrations(ctx context.Context) error {
	migrator := a.diContainer.PGMigrator(ctx)
	if migrator == nil {
		return fmt.Errorf("failed to create migrator")
	}

	err := migrator.Up()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("HTTP server listening on %s", config.AppConfig().HTTP.Address()))

	err := a.apiServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) HTTPServer() *http.Server {
	return a.apiServer
}