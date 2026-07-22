package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	cfg *Config
)

func init() {
	// Load .env if it exists
	_ = godotenv.Load()

	var err error
	cfg, err = LoadConfig()
	if err != nil {
		log.Printf("[INIT] Config error: %v\n", err)
		return
	}

	db, err = sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		log.Printf("[INIT] Database open error: %v\n", err)
		return
	}

	if err := db.Ping(); err != nil {
		log.Printf("[INIT] Database ping error: %v\n", err)
		return
	}

	log.Println("[INIT] Database connected successfully")
	initializeTables()
}

func initializeTables() {
	tables := []string{
		createAccountTableSQL(),
		createTransferTableSQL(),
		createPendingAccountTableSQL(),
		createCouponRedemptionTableSQL(),
		createPendingProfileUpdateTableSQL(),
		createWebAuthnCredentialTableSQL(),
	}

	for _, sql := range tables {
		if _, err := db.Exec(sql); err != nil {
			log.Printf("[INIT] Table creation error: %v\n", err)
		}
	}
}

func createAccountTableSQL() string {
	return `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		number BIGINT,
		email VARCHAR(255),
		encrypted_password VARCHAR(255) NOT NULL,
		balance BIGINT DEFAULT 0,
		has_passkey BOOLEAN DEFAULT false,
		created_at TIMESTAMP
	);`
}

func createTransferTableSQL() string {
	return `CREATE TABLE IF NOT EXISTS transfer (
		id SERIAL PRIMARY KEY,
		transaction_id VARCHAR(32) NOT NULL UNIQUE,
		from_account_id INT NOT NULL REFERENCES account(id),
		to_account_id INT NOT NULL REFERENCES account(id),
		amount BIGINT NOT NULL,
		transaction_type VARCHAR(50) DEFAULT 'transfer',
		status VARCHAR(20) DEFAULT 'completed',
		created_at TIMESTAMP NOT NULL
	);`
}

func createPendingAccountTableSQL() string {
	return `CREATE TABLE IF NOT EXISTS pending_account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		email VARCHAR(255) NOT NULL,
		number BIGINT NOT NULL,
		encrypted_password VARCHAR(255) NOT NULL,
		verification_code VARCHAR(6) NOT NULL,
		expires_at TIMESTAMP NOT NULL
	);`
}

func createCouponRedemptionTableSQL() string {
	return `CREATE TABLE IF NOT EXISTS coupon_redemption (
		id SERIAL PRIMARY KEY,
		coupon_code VARCHAR(255) NOT NULL,
		account_id INT NOT NULL UNIQUE,
		redeemed_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`
}

func createPendingProfileUpdateTableSQL() string {
	return `CREATE TABLE IF NOT EXISTS pending_profile_update (
		id SERIAL PRIMARY KEY,
		account_id INT NOT NULL REFERENCES account(id),
		new_email VARCHAR(255) NOT NULL,
		verification_code VARCHAR(6) NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		UNIQUE (account_id, new_email)
	);`
}

func createWebAuthnCredentialTableSQL() string {
	return `CREATE TABLE IF NOT EXISTS webauthn_credentials (
		id BYTEA PRIMARY KEY,
		public_key BYTEA NOT NULL,
		attestation_type TEXT NOT NULL,
		transport TEXT[] NOT NULL,
		flags JSONB NOT NULL,
		authenticator JSONB NOT NULL,
		account_id INTEGER NOT NULL REFERENCES account(id)
	);`
}

// Handler is the main entry point for Vercel serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS, PUT, PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Check initialization
	if db == nil || cfg == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Server not initialized - check environment variables and database connection",
		})
		return
	}

	// Route the request
	path := strings.TrimPrefix(r.URL.Path, "/api")
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	log.Printf("[REQUEST] %s %s", r.Method, path)

	switch {
	case path == "/health" && r.Method == "GET":
		handleHealth(w, r)

	case path == "/login" && r.Method == "POST":
		handleLogin(w, r)

	case path == "/account" && r.Method == "POST":
		handleCreateAccount(w, r)

	case path == "/account" && r.Method == "GET":
		handleGetCurrentAccount(w, r)

	case path == "/account/verification" && r.Method == "POST":
		handleVerification(w, r)

	case path == "/account/update" && r.Method == "PUT":
		handleUpdateAccount(w, r)

	case path == "/account/transactions" && r.Method == "GET":
		handleGetTransactions(w, r)

	case path == "/account/activity" && r.Method == "GET":
		handleGetActivity(w, r)

	case path == "/transfer" && r.Method == "POST":
		handleTransfer(w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Endpoint not found"})
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Number   int64  `json:"number"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":  "test-token",
		"number": req.Number,
	})
}

func handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Account creation started",
	})
}

func handleGetCurrentAccount(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        1,
		"firstName": "Test",
		"lastName":  "User",
		"email":     "test@example.com",
		"balance":   1000,
	})
}

func handleVerification(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":  "test-token",
		"number": 123456,
	})
}

func handleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Account updated"})
}

func handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]map[string]interface{}{})
}

func handleGetActivity(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]map[string]interface{}{})
}

func handleTransfer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ToAccount int64 `json:"toAccount"`
		Amount    int64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transactionId": "TXN-123456789",
		"amount":        req.Amount,
	})
}

// LoadConfig loads environment variables
type Config struct {
	ListenAddr        string
	DBDSN             string
	JWTSecret         string
	SMTPEmail         string
	SMTPPassword      string
	SMTPHost          string
	SMTPPort          string
	CouponCode        string
	WebAuthnRPOrigin  string
	WebAuthnRPID      string
	WebAuthnDisplayName string
}

func LoadConfig() (*Config, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable must be set")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dbPassword := os.Getenv("DB_PASSWORD")
		if dbPassword == "" {
			return nil, fmt.Errorf("DB_PASSWORD or DATABASE_URL must be set")
		}
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		dbname := getEnv("DB_NAME", "postgres")
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			host, port, user, dbname, dbPassword,
		)
	}

	return &Config{
		ListenAddr:        ":3000",
		DBDSN:             dsn,
		JWTSecret:         jwtSecret,
		SMTPEmail:         os.Getenv("SMTP_EMAIL"),
		SMTPPassword:      os.Getenv("SMTP_PASSWORD"),
		SMTPHost:          getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:          getEnv("SMTP_PORT", "587"),
		CouponCode:        getEnv("COUPON_CODE", "OFFER1000"),
		WebAuthnRPOrigin:  getEnv("WEBAUTHN_RP_ORIGIN", "http://localhost:8080"),
		WebAuthnRPID:      getEnv("WEBAUTHN_RP_ID", "localhost"),
		WebAuthnDisplayName: getEnv("WEBAUTHN_DISPLAY_NAME", "GoBank"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
