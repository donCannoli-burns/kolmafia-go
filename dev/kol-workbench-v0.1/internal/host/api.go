package host

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// API returns a new mux with no mutating routes.
func API() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"status": "ok", "mode": "read-only", "schema_version": 1})
	})
	mux.HandleFunc("GET /api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		home, err := os.UserHomeDir()
		if err != nil {
			http.Error(w, "home unavailable", http.StatusInternalServerError)
			return
		}
		s, err := Collect(r.Context(), home)
		if err != nil {
			http.Error(w, "inventory unavailable", http.StatusInternalServerError)
			return
		}
		writeJSON(w, s)
	})
	return mux
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_ = json.NewEncoder(w).Encode(v)
}

// LoopbackHostGuard blocks DNS-rebinding requests that carry an unexpected Host.
// Future write capabilities will also need authentication and CSRF defenses.
func LoopbackHostGuard(next http.Handler, port string) http.Handler {
	valid := map[string]bool{"127.0.0.1:" + port: true, "localhost:" + port: true}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid[r.Host] || (r.Method != http.MethodGet && r.Method != http.MethodHead) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if origin := r.Header.Get("Origin"); origin != "" {
			allowed := "http://" + r.Host
			if strings.TrimSuffix(origin, "/") != allowed {
				http.Error(w, "forbidden origin", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
