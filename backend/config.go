package main

import (
	"fmt"
	"os"
)

// Config holds all application configuration loaded from environment variables.
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

// getEnv returns the value of the environment variable key, or fallback if not set.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// LoadConfig reads and validates all environment variables, returning a Config
// or an error if any required variable is missing.
func LoadConfig() (*Config, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable must be set")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dbPassword := os.Getenv("DB_PASSWORD")
		if dbPassword == "" {
			return nil, fmt.Errorf("DB_PASSWORD (or DATABASE_URL) environment variable must be set")
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

	// For Vercel deployment, use PORT environment variable if available
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	listenAddr := ":" + port

	return &Config{
		ListenAddr:      listenAddr,
		DBDSN:           dsn,
		JWTSecret:       jwtSecret,
		SMTPEmail:       os.Getenv("SMTP_EMAIL"),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
		SMTPHost:        getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:        getEnv("SMTP_PORT", "587"),
		CouponCode:      getEnv("COUPON_CODE", "OFFER1000"),
		WebAuthnRPOrigin: getEnv("WEBAUTHN_RP_ORIGIN", "http://localhost:8080"),
		WebAuthnRPID:     getEnv("WEBAUTHN_RP_ID", "localhost"),
		WebAuthnDisplayName: getEnv("WEBAUTHN_DISPLAY_NAME", "GoBank"),
	}, nil
}
