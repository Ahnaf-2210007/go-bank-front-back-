package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Global state for the API
var (
	apiServer *APIServer
	store     Storage
	cfg       *Config
)

// Initialize sets up the API server once
func init() {
	// Load .env if it exists (for local development)
	_ = godotenv.Load()

	var err error
	cfg, err = LoadConfig()
	if err != nil {
		log.Printf("Config error: %v\n", err)
		// Don't fatal here, let handler deal with it
	}

	if cfg != nil {
		store, err = NewPostgresStore(cfg.DBDSN)
		if err != nil {
			log.Printf("Storage error: %v\n", err)
		}

		if store != nil {
			apiServer = NewAPIServer(cfg.ListenAddr, store, cfg)
		}
	}
}

// Handler is the main entry point for Vercel serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	// Ensure initialization succeeded
	if apiServer == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Server not initialized - check environment variables",
		})
		return
	}

	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS, PUT, PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.Header().Set("Access-Control-Allow-Credentials", "false")

	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Route the request
	routeRequest(w, r)
}

// routeRequest handles all API routing
func routeRequest(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api")
	method := r.Method

	// Remove trailing slash
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	log.Printf("%s %s", method, path)

	// Route to appropriate handler
	switch {
	case path == "/health" && method == "GET":
		apiServer.handleHealth(w, r)

	case path == "/login" && method == "POST":
		apiServer.handleLogin(w, r)

	case path == "/account" && method == "POST":
		apiServer.handleAccount(w, r)

	case path == "/account" && method == "GET":
		apiServer.handleAccount(w, r)

	case path == "/account/verification" && method == "POST":
		apiServer.handleVerification(w, r)

	case path == "/account/update" && method == "PUT":
		apiServer.handleUpdateAccount(w, r)

	case path == "/account/transactions" && method == "GET":
		apiServer.handleGetTransactions(w, r)

	case path == "/account/activity" && method == "GET":
		apiServer.handleGetActivity(w, r)

	case strings.HasPrefix(path, "/account/") && method == "GET":
		apiServer.handleGetAccountByID(w, r)

	case strings.HasPrefix(path, "/account/") && strings.HasSuffix(path, "/offer") && method == "GET":
		apiServer.handleOffer(w, r)

	case path == "/transfer" && method == "POST":
		apiServer.handleTransfer(w, r)

	case path == "/webauthn/register/begin" && method == "POST":
		apiServer.webAuthnHandler.handleRegisterBegin(w, r)

	case path == "/webauthn/register/finish" && method == "POST":
		apiServer.webAuthnHandler.handleRegisterFinish(w, r)

	case path == "/webauthn/login/begin" && method == "POST":
		apiServer.webAuthnHandler.handleLoginBegin(w, r)

	case strings.HasPrefix(path, "/webauthn/login/finish/") && method == "POST":
		apiServer.webAuthnHandler.handleLoginFinish(w, r)

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
	}
}
