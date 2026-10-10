package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Build variables injected during CI/CD with -ldflags
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
	startTime = time.Now()
)

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageStore struct {
	mu       sync.RWMutex
	messages []Message
	nextID   int
}

var store = &MessageStore{
	messages: []Message{
		{ID: 1, Content: "Hello from Jenkins CI/CD pipeline!", Author: "DevOps Bot", CreatedAt: time.Now().Add(-10 * time.Minute)},
		{ID: 2, Content: "Vue 3 + Go backend deployed successfully.", Author: "Antigravity", CreatedAt: time.Now().Add(-5 * time.Minute)},
	},
	nextID: 3,
}

// corsMiddleware wraps handler to allow frontend requests in dev mode
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"status":    "UP",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(startTime).String(),
	})
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	hostname, _ := os.Hostname()
	jsonResponse(w, http.StatusOK, map[string]any{
		"app":         "Vue + Go Demo Service",
		"version":     Version,
		"git_commit":  GitCommit,
		"build_time":  BuildTime,
		"go_version":  runtime.Version(),
		"os_arch":     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"hostname":    hostname,
		"environment": getEnv("APP_ENV", "production"),
	})
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		store.mu.RLock()
		msgs := make([]Message, len(store.messages))
		copy(msgs, store.messages)
		store.mu.RUnlock()
		jsonResponse(w, http.StatusOK, map[string]any{
			"code": 0,
			"data": msgs,
		})

	case http.MethodPost:
		var req struct {
			Content string `json:"content"`
			Author  string `json:"author"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonResponse(w, http.StatusBadRequest, map[string]any{"error": "Invalid request body"})
			return
		}
		if req.Content == "" {
			jsonResponse(w, http.StatusBadRequest, map[string]any{"error": "Content cannot be empty"})
			return
		}
		if req.Author == "" {
			req.Author = "Anonymous"
		}

		store.mu.Lock()
		newMsg := Message{
			ID:        store.nextID,
			Content:   req.Content,
			Author:    req.Author,
			CreatedAt: time.Now(),
		}
		store.nextID++
		store.messages = append(store.messages, newMsg)
		store.mu.Unlock()

		jsonResponse(w, http.StatusCreated, map[string]any{
			"code": 0,
			"msg":  "created",
			"data": newMsg,
		})

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			jsonResponse(w, http.StatusBadRequest, map[string]any{"error": "id parameter is required"})
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, map[string]any{"error": "invalid id format"})
			return
		}

		store.mu.Lock()
		found := false
		newMsgs := make([]Message, 0, len(store.messages))
		for _, m := range store.messages {
			if m.ID == id {
				found = true
			} else {
				newMsgs = append(newMsgs, m)
			}
		}
		if found {
			store.messages = newMsgs
		}
		store.mu.Unlock()

		if !found {
			jsonResponse(w, http.StatusNotFound, map[string]any{"error": "message not found"})
			return
		}
		jsonResponse(w, http.StatusOK, map[string]any{"code": 0, "msg": "deleted"})

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func setupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/info", infoHandler)
	mux.HandleFunc("/api/messages", messagesHandler)

	return corsMiddleware(mux)
}

func main() {
	port := getEnv("PORT", "8080")
	handler := setupRoutes()

	log.Printf("Starting Vue-Go backend service on port %s...", port)
	log.Printf("Version: %s, GitCommit: %s, BuildTime: %s", Version, GitCommit, BuildTime)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
