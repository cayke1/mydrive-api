package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/cayke1/mydrive-api/internal/config"
	"github.com/cayke1/mydrive-api/internal/folders"
	"github.com/cayke1/mydrive-api/internal/storage"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	log.Printf("Starting MyDrive API on port %s", cfg.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("PostgreSQL ping failed: %v", err)
	}
	log.Println("✓ PostgreSQL connected")

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})
	defer redisClient.Close()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis ping failed: %v", err)
	}
	log.Println("✓ Redis connected")

	storageClient, err := storage.NewMinIOStorage(
		cfg.MinIOEndpoint,
		cfg.MinIOAccessKey,
		cfg.MinIOSecretKey,
		cfg.MinIOBucket,
		cfg.MinIOUseSSL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}
	log.Println("✓ MinIO initialized")

	folderRepo := folders.NewFolderRepository(dbPool)
	folderService := folders.NewFolderService(folderRepo)
	folderController := folders.NewFolderController(folderService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler(dbPool, redisClient, storageClient))
	folderController.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         "127.0.0.1:" + cfg.Port,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	server.Shutdown(shutdownCtx)
	dbPool.Close()
	redisClient.Close()

	log.Println("Server stopped")
}

func healthHandler(db *pgxpool.Pool, redis *redis.Client, st storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		dbOk := db.Ping(ctx) == nil
		redisOk := redis.Ping(ctx).Err() == nil
		storageOk := true

		status := "ok"
		if !dbOk || !redisOk || !storageOk {
			status = "degraded"
		}

		health := map[string]interface{}{
			"status": status,
			"services": map[string]bool{
				"database": dbOk,
				"redis":    redisOk,
				"storage":  storageOk,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if status != "ok" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		json.NewEncoder(w).Encode(health)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
