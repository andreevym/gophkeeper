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
		log.Fatal("init server config", err)
	}

	if _, err := logger.NewLogger(cfg.LogLevel); err != nil {
		log.Fatal("logger can't be initialized:", cfg.LogLevel, err)
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := sqlx.Connect("postgres", cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Failed to connect to database, database uri %s : %v", cfg.DatabaseURI, err)
	}
	defer db.Close()

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Failed to connect to database, database uri %s : %v", cfg.DatabaseURI, err)
	}
	defer conn.Close(ctx)

	err = postgres.Migration(ctx, "migrations", db)
	if err != nil {
		log.Fatalf("Failed to migrate database, database uri %s : %v", cfg.DatabaseURI, err)
	}

	vaultStorage := postgres.NewVaultStorage(db, conn)
	userStorage := postgres.NewUserStorage(db)

	if cfg.JWTSecretKey == "" {
		_, cfg.JWTSecretKey, err = auth.MakeJwtSecretKey()
		if err != nil {
			log.Fatalf("Failed to generate JWT secret key: %v", err)
		}
	}

	jwtPrivateKey, err := auth.ReadJwtSecretKey(cfg.JWTSecretKey)
	if err != nil {
		log.Fatalf("failed to read jwt secret key: %v", err)
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
		defer cancel()
		logger.Logger().Info("listening http server", zap.String("address", cfg.Address))
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("failed listen http server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	for {
		select {
		case <-quit:
			logger.Logger().Info("shutting down server...")
			ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(ctx); err != nil {
				log.Fatalf("Server shutdown failed: %v", err)
			}
			logger.Logger().Info("server http stopped gracefully")
			return
		case <-ctx.Done():
			logger.Logger().Info("shutting down server context done...")
			return
		}
	}
}
