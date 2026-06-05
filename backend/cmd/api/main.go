package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/praxisvr/grama/internal/application/user"
	infraauth "github.com/praxisvr/grama/internal/infrastructure/auth"
	infrahttp "github.com/praxisvr/grama/internal/infrastructure/http"
	"github.com/praxisvr/grama/internal/infrastructure/http/handler"
	"github.com/praxisvr/grama/internal/infrastructure/http/middleware"
	"github.com/praxisvr/grama/internal/infrastructure/persistence/postgres"
)

func main() {
	// Config — fail fast on missing required env vars.
	dbDSN := requireEnv("GRAMA_DB_DSN")
	jwtSecret := requireEnv("JWT_SECRET")
	port := getEnv("SERVER_PORT", "8080")
	allowedOrigins := parseOrigins(getEnv("ALLOWED_ORIGINS", "http://localhost:8100"))

	// Database connection pool.
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}
	log.Println("database connected")

	// Repositories (infrastructure layer).
	userRepo := postgres.NewPostgresUserRepository(pool)
	tokenRepo := postgres.NewPostgresTokenRepository(pool)

	// JWT manager.
	jwtManager := infraauth.NewJWTManager(jwtSecret)

	// Use cases (application layer). Each receives interfaces, never concrete types.
	registerUC := user.NewRegisterUseCase(userRepo)
	loginUC := user.NewLoginUseCase(userRepo, tokenRepo, jwtManager)
	refreshUC := user.NewRefreshUseCase(userRepo, tokenRepo, jwtManager)
	getProfileUC := user.NewGetProfileUseCase(userRepo)
	updateProfileUC := user.NewUpdateProfileUseCase(userRepo)
	createOpUC := user.NewCreateOperatorUseCase(userRepo)
	listOpsUC := user.NewListOperatorsUseCase(userRepo)

	// HTTP layer.
	authMW := middleware.NewAuthMiddleware(jwtManager)
	authHandler := handler.NewAuthHandler(registerUC, loginUC, refreshUC, tokenRepo)
	userHandler := handler.NewUserHandler(getProfileUC, updateProfileUC, createOpUC, listOpsUC)
	router := infrahttp.NewRouter(authHandler, userHandler, authMW, allowedOrigins)

	// HTTP server with sensible timeouts (OWASP recommendation).
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Grama API listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT / SIGTERM.
	quit, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-quit.Done()

	log.Println("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server forced shutdown: %v", err)
	}
	log.Println("server stopped")
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %q is not set", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseOrigins(raw string) []string {
	var origins []string
	for _, o := range splitComma(raw) {
		if o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}

func splitComma(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			parts = append(parts, strings.TrimSpace(s[start:i]))
			start = i + 1
		}
	}
	parts = append(parts, strings.TrimSpace(s[start:]))
	return parts
}
