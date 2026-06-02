package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/praxisvr/grama/internal/infrastructure/config"
)

func main() {
	// Load and validate all environment variables before doing anything else.
	// If any required variable is missing, the server refuses to start.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	// Context tied to OS signals: SIGINT (Ctrl+C) and SIGTERM (Docker stop).
	// When either signal arrives, ctx is cancelled and the server has 10s to finish
	// in-flight requests before we force-close the listener.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Build the connection pool. pgxpool manages multiple connections concurrently —
	// critical for handling simultaneous reservation requests without blocking.
	pool, err := pgxpool.New(ctx, cfg.Database.DSN)
	if err != nil {
		log.Fatalf("cannot create database pool: %v", err)
	}
	defer pool.Close()

	// Verify the DB is reachable before accepting traffic.
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("cannot reach database: %v", err)
	}
	log.Printf("database connected (dsn length: %d chars)", len(cfg.Database.DSN))

	mux := http.NewServeMux()
	registerRoutes(mux, pool)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: mux,
		// Defense: cap request/response timeouts to avoid slow-client attacks (Slowloris).
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start listening in a goroutine so the main goroutine can wait for the signal.
	go func() {
		log.Printf("grama-api listening on :%d [%s]", cfg.Server.Port, cfg.Server.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until SIGINT or SIGTERM.
	<-ctx.Done()
	stop() // release the signal handler

	log.Println("shutdown signal received, draining requests (10s)...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	} else {
		log.Println("server stopped cleanly")
	}
}

func registerRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	mux.HandleFunc("GET /health", healthHandler(pool))
}

// healthHandler returns 200 when the server is up and the DB is reachable,
// or 503 when the DB is unavailable. Used by Docker healthchecks and load balancers.
func healthHandler(pool *pgxpool.Pool) http.HandlerFunc {
	type response struct {
		Status   string `json:"status"`
		Database string `json:"database"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dbStatus := "ok"
		if err := pool.Ping(r.Context()); err != nil {
			dbStatus = "unavailable"
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(response{ //nolint:errcheck
			Status:   "ok",
			Database: dbStatus,
		})
	}
}

