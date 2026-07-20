package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func seedAccount(store Storage, firstName, lastName, email, password string) *Account {
	acc, err := NewAccount(firstName, lastName, email, password)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New account: %d\n", acc.Number)

	return acc
}

func seedAccounts(store Storage) {
	seedAccount(store, "John", "Doe", "john.doe@example.com", "password123")
}

func main() {
	log.Println("Starting GoBank server...")
	
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables.")
	}
	
	// This line is essential. It initializes the random number generator.
	seed := flag.Bool("seed", false, "Seed the database with initial data")
	flag.Parse()

	log.Println("Loading configuration...")
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	log.Printf("Listening on %s", cfg.ListenAddr)

	log.Println("Connecting to database...")
	store, err := NewPostgresStore(cfg)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	log.Println("Database connection established")

	// Initialize the database schema if needed
	log.Println("Creating database tables...")
	tables := []struct {
		name string
		fn   func() error
	}{
		{"accounts", store.CreateAccountTable},
		{"coupon_redemptions", store.CreateCouponRedemptionTable},
		{"pending_accounts", store.CreatePendingAccountTable},
		{"transfers", store.CreateTransferTable},
		{"pending_profile_updates", store.CreatePendingProfileUpdateTable},
		{"webauthn_credentials", store.CreateWebAuthnCredentialTable},
	}

	for _, table := range tables {
		if err := table.fn(); err != nil {
			log.Fatalf("Failed to create %s table: %v", table.name, err)
		}
		log.Printf("Created/verified %s table", table.name)
	}

	//seed stuff
	if *seed {
		fmt.Println("Seeding database with initial data...")
		seedAccounts(store)
	}

	log.Println("Initializing API server...")
	server := NewAPIServer(cfg.ListenAddr, store, cfg)
	log.Println("GoBank server starting...")
	server.Run()
}
