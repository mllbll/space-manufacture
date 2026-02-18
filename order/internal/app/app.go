package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"platform/pkg/closer"
	"platform/pkg/logger"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mllbll/space-manufacture/order/internal/config"
	orderV1 "github.com/mllbll/space-manufacture/shared/pkg/openapi/order/v1"
	"go.uber.org/zap"
)

type App struct {
	diContainer *diContainer

	httpServer *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDi(ctx context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

// Создаем Сервер
func (a *App) initHTTPServer(ctx context.Context) error {
	api := a.diContainer.OrderV1API(ctx)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		logger.Error(ctx, "Ошибка при создании HTTP сервера", zap.Error(err))
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Mount("/", orderServer)

	timeDuration, err := time.ParseDuration(config.AppConfig().OrderHTTTP.ReadHeaderTimeout())
	if err != nil {
		logger.Error(ctx, "Ошибка преобразования времени на чтение сервера в time.time", zap.Error(err))
	}

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().OrderHTTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: timeDuration,
	}

	closer.AddNamed("HTTP Server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}

// Запускаем созданный сервер
func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 http Order server listening on %s", config.AppConfig().OrderHTTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, "❌ Ошибка запуска сервера:", zap.Error(err))
	}

	return nil
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}
