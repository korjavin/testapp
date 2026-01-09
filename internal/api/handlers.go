package api

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/korjavin/testapp/internal/db"
)

//go:embed web/*
var staticFiles embed.FS

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", healthHandler)

	// API routes
	mux.HandleFunc("/api/hello", helloHandler)

	// WebApp routes
	mux.HandleFunc("/", indexHandler)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(); err != nil {
		http.Error(w, fmt.Sprintf(`{"status": "error", "message": %q}`, err.Error()), http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello from testapp!"}`))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serve index.html for all non-file requests
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		data, err := staticFiles.ReadFile("web/index.html")
		if err != nil {
			slog.Error("failed to read index.html", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
		return
	}

	// Try to serve static file
	data, err := staticFiles.ReadFile("web" + r.URL.Path)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// Set content type based on extension
	contentType := getContentType(r.URL.Path)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

func getContentType(path string) string {
	switch {
	case strings.HasSuffix(path, ".css"):
		return "text/css"
	case strings.HasSuffix(path, ".js"):
		return "application/javascript"
	case strings.HasSuffix(path, ".json"):
		return "application/json"
	case strings.HasSuffix(path, ".png"):
		return "image/png"
	case strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(path, ".ico"):
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}

