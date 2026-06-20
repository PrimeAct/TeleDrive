package teldrive

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/tgdrive/teldrive/internal/api/middleware"
	"github.com/tgdrive/teldrive/internal/api/routes"
	"github.com/tgdrive/teldrive/internal/app"
	"github.com/tgdrive/teldrive/internal/config"
	"github.com/tgdrive/teldrive/internal/database"
	"github.com/tgdrive/teldrive/internal/telegram/client"
	"github.com/tgdrive/teldrive/internal/utils/cache"
	"github.com/tgdrive/teldrive/internal/utils/logger"
	"github.com/tgdrive/teldrive/pkg/services"
)

func newRunCommand() *cobra.Command {
	var cfg config.ServerConfig
	loader := config.NewLoader()

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start the TelDrive server",
		Long:  `Start the TelDrive server with the specified configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(cmd.Context(), &cfg)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := loader.Load(cmd, &cfg); err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if err := loader.Validate(&cfg); err != nil {
				return fmt.Errorf("config validation failed: %w", err)
			}
			return nil
		},
	}

	loader.RegisterFlags(cmd.Flags(), reflect.TypeFor[config.ServerConfig]())
	return cmd
}

func runServer(ctx context.Context, cfg *config.ServerConfig) error {
	logLevel, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	log := logger.New(logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
	})
	defer log.Sync()

	log.Info("starting teldrive server",
		"port", cfg.Server.Port,
		"host", cfg.Server.Host,
	)

	db, err := database.Initialize(cfg.DB, log)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer database.Close(db)

	var cacheClient cache.Cache
	if cfg.Redis.Enabled {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Address,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			log.Warn("redis connection failed, falling back to in-memory cache", "error", err)
			cacheClient = cache.NewInMemory()
		} else {
			cacheClient = cache.NewRedis(redisClient)
			defer redisClient.Close()
		}
	} else {
		cacheClient = cache.NewInMemory()
	}

	tgPool, err := client.NewPool(cfg.TG, log)
	if err != nil {
		return fmt.Errorf("failed to initialize telegram client pool: %w", err)
	}
	defer tgPool.Close()

	appCtx := app.NewContext(ctx, db, cacheClient, tgPool, log)
	svc := services.Initialize(appCtx)

	router := chi.NewRouter()
	setupMiddleware(router, cfg, log)
	routes.Register(router, svc, cfg)

	port := cfg.Server.Port
	if cfg.Server.AutoPort {
		port, err = findAvailablePort(port)
		if err != nil {
			return fmt.Errorf("no available port found: %w", err)
		}
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Info("shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("server shutdown error", "error", err)
		}
	}()

	log.Info("server listening", "address", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	wg.Wait()
	return nil
}

func setupMiddleware(r *chi.Mux, cfg *config.ServerConfig, log *zap.Logger) {
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger(log))
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(chimiddleware.Compress(5))

	if cfg.Server.CORS.Enabled {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   cfg.Server.CORS.AllowedOrigins,
			AllowedMethods:   cfg.Server.CORS.AllowedMethods,
			AllowedHeaders:   cfg.Server.CORS.AllowedHeaders,
			ExposedHeaders:   cfg.Server.CORS.ExposedHeaders,
			AllowCredentials: cfg.Server.CORS.AllowCredentials,
			MaxAge:           int(cfg.Server.CORS.MaxAge.Seconds()),
		}))
	}
}

func findAvailablePort(startPort int) (int, error) {
	for port := startPort; port < startPort+1000; port++ {
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			continue
		}
		listener.Close()
		return port, nil
	}
	return 0, fmt.Errorf("no available ports found between %d and %d", startPort, startPort+1000)
}
