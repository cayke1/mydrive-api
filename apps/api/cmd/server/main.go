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

	"github.com/cayke1/mydrive-api/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load()

	app, err := bootstrap()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer app.Close()

	go func() {
		log.Printf("Server listening on %s", app.Server.Addr)
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	app.Server.Shutdown(shutdownCtx)

	log.Println("Server stopped")
}

func healthHandler(db *pgxpool.Pool, redis *redis.Client, st *storage.MinIOStorage) http.HandlerFunc {
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
