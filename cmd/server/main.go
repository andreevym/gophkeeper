package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andreevym/gophkeeper/internal/auth"
	"github.com/andreevym/gophkeeper/internal/config"
	"github.com/andreevym/gophkeeper/internal/handlers"
	"github.com/andreevym/gophkeeper/internal/middleware"
	"github.com/andreevym/gophkeeper/internal/pwd"
	"github.com/andreevym/gophkeeper/internal/server"
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

	s := server.NewServer(router)
	if s == nil {
		log.Fatalf("Server can't be nil: %v", err)
	}
	defer s.Shutdown()
	s.Run(cfg.Address)
}
