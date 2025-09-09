package app

import (
	"context"
	"github.com/ilindan-dev/shortener/internal/config"
	deliveryHTTP "github.com/ilindan-dev/shortener/internal/delivery/http"
	"github.com/ilindan-dev/shortener/internal/logger"
	"github.com/ilindan-dev/shortener/internal/service"
	"github.com/ilindan-dev/shortener/internal/storage/postgres"
	"github.com/ilindan-dev/shortener/internal/storage/redis"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"net/http"
)

// Module is the main Fx module for the Shortener application.
var Module = fx.Options(
	fx.Provide(
		// Core components
		config.NewConfig,
		logger.NewLogger,

		// Storage Layer - concrete implementations
		postgres.NewPool,
		redis.NewClient,

		// Repositories and Caches
		postgres.NewURLRepository,
		postgres.NewClickRepository,
		postgres.NewAnalyticsRepository,
		redis.NewURLCache,

		// Service Layer
		service.NewURLService,
		service.NewAnalyticsService,

		// Delivery Layer
		// We need a special provider for handlers because it needs the baseURL from config.
		func(
			urlService *service.URLService,
			analyticsService *service.AnalyticsService,
			logger *zerolog.Logger,
			cfg *config.Config,
		) *deliveryHTTP.Handlers {
			return deliveryHTTP.NewHandlers(urlService, analyticsService, logger, cfg.HTTP.BaseURL)
		},
		deliveryHTTP.NewServer,
	),
	// This invoke bootstraps the HTTP server.
	fx.Invoke(func(server *deliveryHTTP.Server, lc fx.Lifecycle) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						panic(err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
		})
	}),
)
