package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/rayfiyo/go-layered-architecture/internal/application/service"
	"github.com/rayfiyo/go-layered-architecture/internal/infrastructure/repository"
	"github.com/rayfiyo/go-layered-architecture/internal/presentation/handler"
)

func main() {
	// ---- Config ----
	addr := envOrDefault("APP_ADDR", ":8080")
	dsn := envOrDefault(
		"DB_DSN", "file:app.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")

	// ---- DB ----
	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("failed to close db: %v", cerr)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	if err := ensureSchema(db); err != nil {
		log.Fatalf("failed to ensure schema: %v", err)
	}

	// ---- Wiring (DIP) ----
	userRepo := repository.NewUserRepositorySQLX(db)
	userSvc := service.NewUserService(userRepo)
	r := handler.NewRouter(userSvc)

	// ---- HTTP server ----
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// ---- Graceful shutdown ----
	ctx, stop := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(
			err, http.ErrServerClosed,
		) {
			log.Printf("server error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Printf("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	log.Printf("bye")
}

func envOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func ensureSchema(db *sqlx.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
`
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("exec schema: %w", err)
	}
	return nil
}
