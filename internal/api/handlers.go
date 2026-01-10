package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/korjavin/testapp/internal/db"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", healthHandler)

	// API routes
	mux.HandleFunc("/api/hello", helloHandler)

	// WebApp routes - serve static files from web/ directory
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/css/", staticHandler)
	mux.HandleFunc("/js/", staticHandler)

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
	// Serve index.html for root path
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		http.ServeFile(w, r, "web/index.html")
		return
	}

	// Let staticHandler handle other files
	staticHandler(w, r)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	// Remove prefix and serve from web/ directory
	path := strings.TrimPrefix(r.URL.Path, "/")
	if strings.HasPrefix(path, "css/") {
		path = strings.TrimPrefix(path, "css/")
		http.ServeFile(w, r, "web/css/"+path)
		return
	}
	if strings.HasPrefix(path, "js/") {
		path = strings.TrimPrefix(path, "js/")
		http.ServeFile(w, r, "web/js/"+path)
		return
	}

	// Default: serve from web/
	http.ServeFile(w, r, "web/"+path)
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
