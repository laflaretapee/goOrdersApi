package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/laflaretapee/go-orders-api/internal/config"
	"github.com/laflaretapee/go-orders-api/internal/httpapi"
	"github.com/laflaretapee/go-orders-api/internal/order"
	"github.com/laflaretapee/go-orders-api/internal/storage/postgres"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := db.PingContext(pingCtx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	repository := postgres.NewOrderRepository(db)
	service := order.NewService(repository)
	handler := httpapi.NewHandler(service)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("starting API on http://localhost:%s", cfg.Port)
		log.Printf("database: %s", redactDatabaseURL(cfg.DatabaseURL))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
}

func redactDatabaseURL(raw string) string {
	if raw == "" {
		return ""
	}

	if idx := len("postgres://"); len(raw) > idx && raw[:idx] == "postgres://" {
		return "postgres://***"
	}

	if idx := len("postgresql://"); len(raw) > idx && raw[:idx] == "postgresql://" {
		return "postgresql://***"
	}

	if os.Getenv("DATABASE_URL") != "" {
		return "DATABASE_URL is set"
	}

	return "configured"
}
