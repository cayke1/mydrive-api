package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/cayke1/mydrive-api/internal/auth"
	"github.com/cayke1/mydrive-api/internal/config"
	"github.com/cayke1/mydrive-api/internal/files"
	"github.com/cayke1/mydrive-api/internal/folders"
	"github.com/cayke1/mydrive-api/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	DB      *pgxpool.Pool
	Redis   *redis.Client
	Storage *storage.MinIOStorage
	Config  config.Config
	Mux     *http.ServeMux
	Server  *http.Server
}

func bootstrap() (*App, error) {
	cfg := config.Load()

	log.Printf("Starting MyDrive API on port %s", cfg.Port)

	app := &App{Config: cfg}

	if err := app.initDatabase(); err != nil {
		return nil, err
	}

	if err := app.initRedis(); err != nil {
		return nil, err
	}

	if err := app.initStorage(); err != nil {
		return nil, err
	}

	app.registerRoutes()
	app.createServer()

	return app, nil
}

func (a *App) initDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, a.Config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to create PostgreSQL pool: %v", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("PostgreSQL ping failed: %v", err)
	}

	a.DB = dbPool
	log.Println("✓ PostgreSQL connected")
	return nil
}

func (a *App) initRedis() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	redisClient := redis.NewClient(&redis.Options{
		Addr: a.Config.RedisURL,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis ping failed: %v", err)
	}

	a.Redis = redisClient
	log.Println("✓ Redis connected")
	return nil
}

func (a *App) initStorage() error {
	storageClient, err := storage.NewMinIOStorage(
		a.Config.MinIOEndpoint,
		a.Config.MinIOAccessKey,
		a.Config.MinIOSecretKey,
		a.Config.MinIOBucket,
		a.Config.MinIOUseSSL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	a.Storage = storageClient
	log.Println("✓ MinIO initialized")
	return nil
}

func (a *App) registerRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler(a.DB, a.Redis, a.Storage))

	folderRepo := folders.NewFolderRepository(a.DB)
	folderService := folders.NewFolderService(folderRepo)
	folderController := folders.NewFolderController(folderService)

	fileRepo := files.NewFileRepository(a.DB)
	fileService := files.NewFileService(fileRepo)
	fileController := files.NewFileController(fileService)

	usersRepo := auth.NewUserRepository(a.DB)
	authService := auth.NewAuthService(usersRepo)
	authController := auth.NewAuthController(authService)

	folderController.RegisterRoutes(mux)
	fileController.RegisterRoutes(mux)
	authController.RegisterRoutes(mux)

	a.Mux = mux
}

func (a *App) createServer() {
	a.Server = &http.Server{
		Addr:         "127.0.0.1:" + a.Config.Port,
		Handler:      corsMiddleware(a.Mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
	if a.Redis != nil {
		a.Redis.Close()
	}
}
