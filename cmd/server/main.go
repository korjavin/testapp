package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/korjavin/testapp/internal/api"
	"github.com/korjavin/testapp/internal/db"
	"github.com/korjavin/testapp/internal/middleware"
)

func main() {
	// Flags
	migrate := flag.Bool("migrate", false, "Run database migrations only")
	flag.Parse()

	// Setup logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	// Load environment
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data.db"
	}

	if err := db.InitDatabase(dbPath); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations only if flag is set
	if *migrate {
		slog.Info("Migrations completed")
		return
	}

	// Setup routes
	mux := api.SetupRoutes()

	// Apply middleware
	handler := middleware.Health(mux)
	handler = middleware.Recovery(handler)
	handler = middleware.Logger(handler)

	// Create server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("Server starting", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server stopped")
	fmt.Println("Server stopped")
}
