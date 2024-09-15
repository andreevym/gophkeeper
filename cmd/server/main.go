package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andreevym/gophkeeper/internal/auth"
	"github.com/andreevym/gophkeeper/internal/config"
	"github.com/andreevym/gophkeeper/internal/handlers"
	"github.com/andreevym/gophkeeper/internal/middleware"
	"github.com/andreevym/gophkeeper/internal/pwd"
	"github.com/andreevym/gophkeeper/internal/storage/postgres"
	"github.com/andreevym/gophkeeper/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var buildVersion string
var buildDate string
var buildCommit string

func printVersion() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func main() {
	printVersion()

	cfg, err := config.NewServerConfig().Init()
	if err != nil {
		log.Fatalf("Error initializing server config: %v", err)
	}

	_, err = logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Error initializing logger with level %s: %v", cfg.LogLevel, err)
	}

	db, err := sqlx.Connect("postgres", cfg.DatabaseURI)
	if err != nil {
		logger.Logger().Fatal("Failed to connect to database", zap.String("databaseURI", cfg.DatabaseURI), zap.Error(err))
	}
	defer db.Close()

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		logger.Logger().Fatal("Failed to connect to database", zap.String("databaseURI", cfg.DatabaseURI), zap.Error(err))
	}
	defer conn.Close(ctx)

	err = postgres.Migration(ctx, "migrations", db)
	if err != nil {
		logger.Logger().Fatal("Failed to migrate database", zap.String("databaseURI", cfg.DatabaseURI), zap.Error(err))
	}

	vaultStorage := postgres.NewVaultStorage(db, conn)
	userStorage := postgres.NewUserStorage(db)

	if cfg.JWTSecretKey == "" {
		_, cfg.JWTSecretKey, err = auth.MakeJwtSecretKey()
		if err != nil {
			logger.Logger().Fatal("Failed to generate JWT secret key", zap.Error(err))
		}
	}

	jwtPrivateKey, err := auth.ReadJwtSecretKey(cfg.JWTSecretKey)
	if err != nil {
		logger.Logger().Fatal("Failed to read JWT secret key", zap.Error(err))
	}

	authProvider := auth.NewAuthProvider(userStorage, jwtPrivateKey)
	authMiddleware := auth.NewAuthMiddleware(authProvider, cfg.JWTSecretKey, handlers.AuthSignInURI, handlers.AuthSignUpURI)
	hashService := pwd.NewHashService()
	serviceHandlers := handlers.NewServiceHandlers(db, authProvider, vaultStorage, userStorage, hashService)

	router := handlers.NewRouter(
		serviceHandlers,
		authMiddleware.WithAuthentication,
		middleware.WithRequestLoggerMiddleware,
	)

	logger.Logger().Info("Server listening", zap.String("addr", cfg.Address))

	httpServer := &http.Server{Addr: cfg.Address, Handler: router}
	go func() {
		defer logger.Logger().Info("HTTP server stopped gracefully")
		logger.Logger().Info("Listening HTTP server", zap.String("address", cfg.Address))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger().Fatal("HTTP server listen failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	for {
		select {
		case <-quit:
			logger.Logger().Info("Shutting down server...")
			ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(ctx); err != nil {
				logger.Logger().Fatal("Server shutdown failed", zap.Error(err))
			}
			return
		case <-ctx.Done():
			logger.Logger().Info("Server context done")
			return
		}
	}
}
