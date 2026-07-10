package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/mail"
	"net/smtp"
	"sort"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type APIServer struct {
	listenAddr      string
	store           Storage
	cfg             *Config
	webAuthnHandler *WebAuthnHandler
}

func NewAPIServer(listenAddr string, store Storage, cfg *Config) *APIServer {
	webAuthnHandler, err := NewWebAuthnHandler(store, cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &APIServer{
		listenAddr:      listenAddr,
		store:           store,
		cfg:             cfg,
		webAuthnHandler: webAuthnHandler,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", makeHTTPHandlerFunc(s.handleHealth))

	// Authentication routes
	router.HandleFunc("/login", makeHTTPHandlerFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/verification", makeHTTPHandlerFunc(s.handleVerification))
	router.HandleFunc("/account/update", makeHTTPHandlerFunc(s.handleUpdateAccount))
	router.HandleFunc("/account/transactions", makeHTTPHandlerFunc(s.handleGetTransactions))
	router.HandleFunc("/account/activity", makeHTTPHandlerFunc(s.handleGetActivity))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandlerFunc(s.handleGetAccountByID), s.store, s.cfg.JWTSecret))
	router.HandleFunc("/account/{id}/offer", withJWTAuth(makeHTTPHandlerFunc(s.handleOffer), s.store, s.cfg.JWTSecret))
	router.HandleFunc("/transfer", makeHTTPHandlerFunc(s.handleTransfer))

	// WebAuthn routes
	router.HandleFunc("/webauthn/register/begin", makeHTTPHandlerFunc(s.webAuthnHandler.handleRegisterBegin))
	router.HandleFunc("/webauthn/register/finish", makeHTTPHandlerFunc(s.webAuthnHandler.handleRegisterFinish))
	router.HandleFunc("/webauthn/login/begin", makeHTTPHandlerFunc(s.webAuthnHandler.handleLoginBegin))
	router.HandleFunc("/webauthn/login/finish/{email}", makeHTTPHandlerFunc(s.webAuthnHandler.handleLoginFinish))

	// Apply middleware
	handler := s.corsMiddleware(router)

	log.Println("JSON API server is running on", s.listenAddr)

	if err := http.ListenAndServe(s.listenAddr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// corsMiddleware adds CORS headers to all responses and handles preflight requests
func (s *APIServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS, PUT, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Credentials", "false")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleHealth returns the health status of the API server
func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	// Handle login API requests here
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.store.GetAccountByNumber(req.Number)
	if err != nil {
		return fmt.Errorf("not authorized")
	}

	if !acc.ValidPassword(req.Password) {
		return fmt.Errorf("not authorized")
	}

	token, err := createJWT(acc, s.cfg.JWTSecret)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Token:  token,
		Number: acc.Number,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// Handle account-related API requests here
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "GET" {
		return s.handleGetCurrentAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetCurrentAccount(w http.ResponseWriter, r *http.Request) error {
	acc := s.getAuthenticatedAccount(w, r)
	if acc == nil {
		return nil
	}

	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	// Create the specific account from the image, instead of a random one.
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// Handle account creation API requests here
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	if req.FirstName == "" || req.LastName == "" {
		return fmt.Errorf("first name and last name are required")
	}
	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf("a valid email address is required")
	}

	account, err := NewAccount(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return err
	}

	code, err := generateVerificationCode()
	if err != nil {
		return err
	}

	pending := &PendingAccount{
		FirstName:         account.FirstName,
		LastName:          account.LastName,
		Email:             account.Email,
		Number:            account.Number,
		EncryptedPassword: account.EncryptedPassword,
		VerificationCode:  code,
		ExpiresAt:         time.Now().UTC().Add(5 * time.Minute),
	}

	if err := s.store.CreatePendingAccount(pending); err != nil {
		return err
	}

	if err := sendVerificationEmail(req.Email, code, s.cfg); err != nil {
		// Roll back the pending account so the user can retry.
		_ = s.store.DeletePendingAccount(pending.ID)
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message": "verification code sent to email",
	})
}

// handleVerification verifies the 6-digit code sent to the user's email and,
// if valid, creates the account and returns it.
func (s *APIServer) handleVerification(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req VerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	if len(req.Code) != 6 {
		return fmt.Errorf("verification code must be 6 digits")
	}

	pending, err := s.store.GetPendingAccountByCode(req.Code)
	if err != nil {
		return err
	}

	if time.Now().UTC().After(pending.ExpiresAt) {
		_ = s.store.DeletePendingAccount(pending.ID)
		return fmt.Errorf("verification code has expired")
	}

	account := &Account{
		FirstName:         pending.FirstName,
		LastName:          pending.LastName,
		Email:             pending.Email,
		Number:            pending.Number,
		EncryptedPassword: pending.EncryptedPassword,
		CreatedAt:         time.Now().UTC(),
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	if err := s.store.DeletePendingAccount(pending.ID); err != nil {
		log.Printf("failed to delete pending account %d: %v", pending.ID, err)
	}

	// Send signup confirmation email with account details
	go func(acc *Account, cfg *Config) {
		if err := sendSignupConfirmationEmail(acc, cfg); err != nil {
			log.Printf("failed to send signup confirmation email to %s: %v", acc.Email, err)
		}
	}(account, s.cfg)

	token, err := createJWT(account, s.cfg.JWTSecret)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, LoginResponse{
		Token:  token,
		Number: account.Number,
	})
}

// generateVerificationCode returns a cryptographically random 6-digit string (zero-padded).
func generateVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// sendVerificationEmail sends the verification code to the given email address
// via Gmail SMTP (smtp.gmail.com:587). The sender credentials are read from
// SMTP_EMAIL and SMTP_PASSWORD. When either variable is not set the code is
// only logged so the server still starts in development environments.
func sendVerificationEmail(email, code string, cfg *Config) error {
	if cfg.SMTPEmail == "" || cfg.SMTPPassword == "" {
		log.Printf("SMTP_EMAIL or SMTP_PASSWORD not set - verification code for %s: %s", email, code)
		return nil
	}

	host := cfg.SMTPHost
	port := cfg.SMTPPort
	addr := host + ":" + port
	auth := smtp.PlainAuth("", cfg.SMTPEmail, cfg.SMTPPassword, host)

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }
        .container { max-width: 600px; margin: 0 auto; background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #4ea2ff 0%%, #0066cc 100%%); padding: 40px 20px; text-align: center; }
        .logo { font-size: 28px; font-weight: bold; color: white; margin-bottom: 10px; }
        .tagline { color: rgba(255,255,255,0.9); font-size: 14px; }
        .content { padding: 40px 30px; }
        .greeting { font-size: 24px; font-weight: 600; color: #1a1a1a; margin-bottom: 15px; }
        .description { font-size: 15px; color: #666; line-height: 1.6; margin-bottom: 30px; }
        .code-box { background: linear-gradient(135deg, #4ea2ff15 0%%, #0066cc10 100%%); border: 1px solid #4ea2ff30; border-radius: 8px; padding: 20px; text-align: center; margin: 30px 0; }
        .code-label { font-size: 12px; color: #666; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 10px; }
        .code { font-size: 36px; font-weight: 700; color: #4ea2ff; letter-spacing: 4px; font-family: 'Courier New', monospace; }
        .info-box { background: #f9f9f9; border-left: 4px solid #4ea2ff; padding: 15px; margin: 20px 0; border-radius: 4px; }
        .info-box p { margin: 0; font-size: 14px; color: #666; }
        .footer { background: #f5f5f5; padding: 20px 30px; border-top: 1px solid #eee; font-size: 12px; color: #999; text-align: center; }
        .button { display: inline-block; background: linear-gradient(135deg, #4ea2ff 0%%, #0066cc 100%%); color: white; padding: 12px 30px; border-radius: 6px; text-decoration: none; font-weight: 600; margin-top: 10px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">🏦 GoBank</div>
            <div class="tagline">Secure Banking Made Simple</div>
        </div>
        <div class="content">
            <p class="greeting">Verify Your Email</p>
            <p class="description">
                We're excited to have you join GoBank. To complete your signup and secure your account, please use the verification code below:
            </p>
            <div class="code-box">
                <div class="code-label">Your Verification Code</div>
                <div class="code">%s</div>
            </div>
            <div class="info-box">
                <p><strong>⏰ Code Expires:</strong> This code is valid for 5 minutes. If you don't use it in time, request a new one.</p>
            </div>
            <p style="font-size: 14px; color: #666; margin-top: 20px;">
                <strong>Didn't sign up?</strong><br>
                If you didn't create this account, you can safely ignore this email. Your account won't be created until you verify your email.
            </p>
        </div>
        <div class="footer">
            <p>© 2024 GoBank. All rights reserved. | <a href="#" style="color: #4ea2ff; text-decoration: none;">Privacy Policy</a> | <a href="#" style="color: #4ea2ff; text-decoration: none;">Terms of Service</a></p>
        </div>
    </div>
</body>
</html>`, code)

	msg := strings.Join([]string{
		"From: " + cfg.SMTPEmail,
		"To: " + email,
		"Subject: Your GoBank Verification Code",
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
		"",
		htmlBody,
	}, "\r\n")

	if err := smtp.SendMail(addr, auth, cfg.SMTPEmail, []string{email}, []byte(msg)); err != nil {
		return err
	}
	log.Printf("verification email sent to %s via SMTP", email)
	return nil
}

func generateTransactionID() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("TXN-%012d", n.Int64()), nil
}

func sendTransferNotificationEmail(email string, isRecipient bool, recipientName, senderName string, transfer *TransferResult, recipient *Account, cfg *Config) error {
	if strings.TrimSpace(email) == "" {
		log.Printf("transfer notification skipped for transaction %s: missing recipient email", transfer.TransactionID)
		return nil
	}

	if cfg.SMTPEmail == "" || cfg.SMTPPassword == "" {
		log.Printf(
			"SMTP_EMAIL or SMTP_PASSWORD not set - transfer notification for %s: transaction=%s amount=%d at=%s",
			email,
			transfer.TransactionID,
			transfer.Amount,
			transfer.CreatedAt.Format(time.RFC3339),
		)
		return nil
	}

	var subject, greeting, actionText string
	var amountColor string = "#059669"

	if isRecipient {
		subject = fmt.Sprintf("Payment Received - $%.2f from %s", float64(transfer.Amount)/100.0, senderName)
		greeting = fmt.Sprintf("Hello %s,", recipientName)
		actionText = fmt.Sprintf("You have received a transfer of <strong style=\"color: %s;\">$%.2f</strong> from %s.", amountColor, float64(transfer.Amount)/100.0, senderName)
	} else {
		subject = fmt.Sprintf("Transfer Complete - $%.2f sent to %s", float64(transfer.Amount)/100.0, recipientName)
		greeting = fmt.Sprintf("Hello %s,", senderName)
		actionText = fmt.Sprintf("You have successfully transferred <strong style=\"color: %s;\">$%.2f</strong> to %s.", amountColor, float64(transfer.Amount)/100.0, recipientName)
	}

	formattedAmount := fmt.Sprintf("$%.2f", float64(transfer.Amount)/100.0)
	formattedDate := transfer.CreatedAt.Format("January 2, 2006 at 3:04 PM MST")

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', sans-serif; background: #f5f5f5; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #4ea2ff 0%%, #0066cc 100%%); padding: 40px 20px; text-align: center; color: white; }
        .logo { font-size: 24px; font-weight: bold; margin-bottom: 8px; }
        .tagline { font-size: 13px; opacity: 0.9; }
        .content { padding: 40px 30px; }
        .greeting { font-size: 18px; font-weight: 600; color: #1a1a1a; margin-bottom: 20px; }
        .action-text { font-size: 15px; color: #333; line-height: 1.6; margin-bottom: 30px; }
        .transaction-box { background: #f9fafb; border-left: 4px solid #4ea2ff; padding: 20px; border-radius: 6px; margin: 20px 0; }
        .transaction-item { display: flex; justify-content: space-between; margin: 12px 0; font-size: 14px; }
        .transaction-label { color: #666; font-weight: 500; }
        .transaction-value { color: #1a1a1a; font-weight: 600; }
        .amount-display { font-size: 28px; font-weight: 700; color: #059669; margin: 15px 0; text-align: center; }
        .reference-box { background: #f0f4ff; border: 1px solid #dde7ff; border-radius: 6px; padding: 12px; margin: 15px 0; }
        .reference-label { font-size: 12px; color: #666; text-transform: uppercase; letter-spacing: 0.5px; }
        .reference-value { font-size: 13px; color: #1a1a1a; font-family: 'Courier New', monospace; font-weight: 600; margin-top: 4px; word-break: break-all; }
        .info-box { background: #fffbf0; border-left: 4px solid #ff9800; padding: 15px; margin: 20px 0; border-radius: 4px; }
        .info-box p { margin: 8px 0; font-size: 13px; color: #555; }
        .info-box strong { color: #ff9800; }
        .cta-button { display: inline-block; background: #4ea2ff; color: white; padding: 12px 30px; border-radius: 6px; text-decoration: none; font-weight: 600; font-size: 14px; margin: 20px 0; text-align: center; }
        .footer { background: #f5f5f5; padding: 20px 30px; border-top: 1px solid #eee; font-size: 12px; color: #999; text-align: center; }
        .footer a { color: #4ea2ff; text-decoration: none; }
    </style>
</head>
<body>
<div class="container">
    <div class="header">
        <div class="logo">GoBank</div>
        <div class="tagline">Secure Banking Services</div>
    </div>
    
    <div class="content">
        <p class="greeting">%s</p>
        
        <p class="action-text">%s</p>
        
        <div class="transaction-box">
            <div class="transaction-item">
                <span class="transaction-label">Transaction ID:</span>
                <span class="transaction-value">%s</span>
            </div>
            <div class="transaction-item">
                <span class="transaction-label">Amount:</span>
                <span class="transaction-value">%s</span>
            </div>
            <div class="transaction-item">
                <span class="transaction-label">Date & Time:</span>
                <span class="transaction-value">%s</span>
            </div>
        </div>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="#" class="cta-button">View Account Details</a>
        </div>
        
        <div class="info-box">
            <p><strong>🔒 Security Note:</strong> This is an automated notification from GoBank. Please do not reply to this email. For security inquiries, log in to your account or contact our support team.</p>
        </div>
        
        <p style="font-size: 13px; color: #666; margin-top: 20px;">
            <strong>Questions?</strong><br>
            If you don't recognize this transaction or have concerns, please <a href="#" style="color: #4ea2ff; text-decoration: none;">contact our support team</a> immediately.
        </p>
    </div>
    
    <div class="footer">
        <p>© 2024 GoBank. All rights reserved.</p>
        <p><a href="#">Privacy Policy</a> | <a href="#">Terms of Service</a> | <a href="#">Contact Us</a></p>
    </div>
</div>
</body>
</html>`, greeting, actionText, transfer.TransactionID, formattedAmount, formattedDate)

	msg := strings.Join([]string{
		"From: " + cfg.SMTPEmail,
		"To: " + email,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
		"",
	}, "\r\n")
	msg += "\r\n" + htmlBody

	host := cfg.SMTPHost
	port := cfg.SMTPPort
	addr := host + ":" + port
	auth := smtp.PlainAuth("", cfg.SMTPEmail, cfg.SMTPPassword, host)

	if err := smtp.SendMail(addr, auth, cfg.SMTPEmail, []string{email}, []byte(msg)); err != nil {
		return err
	}

	log.Printf("transfer notification sent to %s for transaction %s", email, transfer.TransactionID)
	return nil
}

func maskedAccountNumber(number int64) string {
	numberStr := strconv.FormatInt(number, 10)
	if len(numberStr) <= 4 {
		return numberStr
	}
	return strings.Repeat("*", len(numberStr)-4) + numberStr[len(numberStr)-4:]
}

// sendSignupConfirmationEmail sends a welcome email after successful account verification
// with account details including name, account number, ID, and current balance
func sendSignupConfirmationEmail(account *Account, cfg *Config) error {
	if cfg.SMTPEmail == "" || cfg.SMTPPassword == "" {
		log.Printf("SMTP_EMAIL or SMTP_PASSWORD not set - signup confirmation for %s: number=%d id=%d balance=%d", account.Email, account.Number, account.ID, account.Balance)
		return nil
	}

	host := cfg.SMTPHost
	port := cfg.SMTPPort
	addr := host + ":" + port
	auth := smtp.PlainAuth("", cfg.SMTPEmail, cfg.SMTPPassword, host)

	fullName := strings.TrimSpace(account.FirstName + " " + account.LastName)
	if fullName == "" {
		fullName = "GoBank User"
	}

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }
        .container { max-width: 600px; margin: 0 auto; background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #4ea2ff 0%%, #0066cc 100%%); padding: 40px 20px; text-align: center; }
        .logo { font-size: 28px; font-weight: bold; color: white; margin-bottom: 10px; }
        .tagline { color: rgba(255,255,255,0.9); font-size: 14px; }
        .content { padding: 40px 30px; }
        .greeting { font-size: 24px; font-weight: 600; color: #1a1a1a; margin-bottom: 15px; }
        .description { font-size: 15px; color: #666; line-height: 1.6; margin-bottom: 20px; }
        .section-title { font-size: 14px; font-weight: 700; color: #4ea2ff; text-transform: uppercase; letter-spacing: 1px; margin: 25px 0 15px 0; }
        .detail-item { background: #f9f9f9; border-left: 4px solid #4ea2ff; padding: 12px 15px; margin: 10px 0; border-radius: 4px; }
        .detail-label { font-size: 12px; color: #999; text-transform: uppercase; letter-spacing: 0.5px; }
        .detail-value { font-size: 16px; font-weight: 600; color: #1a1a1a; margin-top: 4px; }
        .feature-box { background: #f0f4ff; border: 1px solid #dde7ff; border-radius: 8px; padding: 15px; margin: 10px 0; }
        .feature-box strong { color: #4ea2ff; }
        .feature-box p { margin: 5px 0; font-size: 13px; color: #555; }
        .security-box { background: #fffbf0; border-left: 4px solid #ff9800; padding: 15px; margin: 20px 0; border-radius: 4px; }
        .security-box strong { color: #ff9800; }
        .security-box p { margin: 5px 0; font-size: 13px; color: #555; }
        .footer { background: #f5f5f5; padding: 20px 30px; border-top: 1px solid #eee; font-size: 12px; color: #999; text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">🏦 GoBank</div>
            <div class="tagline">Welcome to Secure Banking</div>
        </div>
        <div class="content">
            <p class="greeting">Welcome, %s! 🎉</p>
            <p class="description">
                Your GoBank account has been successfully created and verified. You're all set to start managing your finances securely.
            </p>
            
            <div class="section-title">📋 Your Account Details</div>
            <div class="detail-item">
                <div class="detail-label">Account Holder</div>
                <div class="detail-value">%s</div>
            </div>
            <div class="detail-item">
                <div class="detail-label">Account Number</div>
                <div class="detail-value">%d</div>
            </div>
            <div class="detail-item">
                <div class="detail-label">Account ID</div>
                <div class="detail-value">#%d</div>
            </div>
            <div class="detail-item">
                <div class="detail-label">Current Balance</div>
                <div class="detail-value">$%.2f</div>
            </div>

            <div class="section-title">🚀 Next Steps</div>
            <div class="feature-box">
                <strong>1. Complete Your Profile</strong>
                <p>Add a profile picture and additional information to personalize your account.</p>
            </div>
            <div class="feature-box">
                <strong>2. Register a Passkey</strong>
                <p>Set up a passkey for faster and more secure login. No more passwords needed!</p>
            </div>
            <div class="feature-box">
                <strong>3. Explore Features</strong>
                <p>Start making transfers, check your history, and discover exclusive GoBank offers.</p>
            </div>

            <div class="security-box">
                <strong>🔒 Security Reminder</strong>
                <p>• Never share your password or verification codes with anyone</p>
                <p>• Always enable two-factor authentication when available</p>
                <p>• Only access GoBank through official channels</p>
                <p>• Logout when using shared computers</p>
            </div>

            <p style="font-size: 13px; color: #666; margin-top: 20px; text-align: center;">
                <strong>Need Help?</strong><br>
                Visit our support center or contact us at support@gobank.com
            </p>
        </div>
        <div class="footer">
            <p>© 2024 GoBank. All rights reserved. | <a href="#" style="color: #4ea2ff; text-decoration: none;">Privacy Policy</a> | <a href="#" style="color: #4ea2ff; text-decoration: none;">Terms of Service</a></p>
        </div>
    </div>
</body>
</html>`, fullName, fullName, account.Number, account.ID, float64(account.Balance)/100)

	msg := strings.Join([]string{
		"From: " + cfg.SMTPEmail,
		"To: " + account.Email,
		"Subject: Welcome to GoBank - Account Confirmation",
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
		"",
		htmlBody,
	}, "\r\n")

	if err := smtp.SendMail(addr, auth, cfg.SMTPEmail, []string{account.Email}, []byte(msg)); err != nil {
		return err
	}

	log.Printf("signup confirmation email sent to %s for account %d", account.Email, account.ID)
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	// Handle account deletion API requests here
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleOffer(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	id, err := getID(r)
	if err != nil {
		return err
	}

	var req OfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	if req.CouponCode == "" {
		return fmt.Errorf("coupon code is required")
	}

	if req.CouponCode != s.cfg.CouponCode {
		return fmt.Errorf("invalid coupon code")
	}

	if err := s.store.RedeemCoupon(id, req.CouponCode); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "offer applied, 1000 added to your balance"})
}

// handleGetTransactions handles GET /account/transactions.
// It returns the paginated transaction history for the authenticated user.
// Query parameters:
//   - limit  (default 20, max 100): number of records per page
//   - offset (default 0): number of records to skip
//   - type   (optional): filter by transaction type (e.g. "transfer")
//   - month  (optional): filter by month, as a number (1-12) or name (e.g. "march")
func (s *APIServer) handleGetTransactions(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	tokenString := r.Header.Get("Authorization")
	token, err := validateJWT(tokenString, s.cfg.JWTSecret)
	if err != nil || !token.Valid {
		WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "not authorized"})
		return nil
	}

	claims := token.Claims.(jwt.MapClaims)
	accountNumber := int64(claims["accountNumber"].(float64))

	acc, err := s.store.GetAccountByNumber(accountNumber)
	if err != nil {
		return fmt.Errorf("account not found")
	}

	query := r.URL.Query()

	limit := 20
	if l := query.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if limit > 100 {
		limit = 100
	}

	offset := 0
	if o := query.Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	txType := query.Get("type")

	month := 0
	if m := query.Get("month"); m != "" {
		month = parseMonth(m)
	}

	transactions, err := s.store.GetTransactionHistory(acc.ID, limit, offset, txType, month)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, transactions)
}

// ActivityEvent represents a single event in the user's activity timeline.
// It includes transactions, profile updates, and passkey registrations.
type ActivityEvent struct {
	Type      string    `json:"type"`      // "transfer_sent", "transfer_received", "profile_update", "passkey_registered", "account_created"
	Title     string    `json:"title"`
	Details   string    `json:"details"`
	Amount    int64     `json:"amount,omitempty"`    // For transfers only
	Status    string    `json:"status"`    // "completed", "pending", "synced"
	Timestamp time.Time `json:"timestamp"`
}

// handleGetActivity handles GET /account/activity.
// It returns recent activity events for the authenticated user, including transfers and profile changes.
func (s *APIServer) handleGetActivity(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	tokenString := r.Header.Get("Authorization")
	token, err := validateJWT(tokenString, s.cfg.JWTSecret)
	if err != nil || !token.Valid {
		WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "not authorized"})
		return nil
	}

	claims := token.Claims.(jwt.MapClaims)
	accountNumber := int64(claims["accountNumber"].(float64))

	acc, err := s.store.GetAccountByNumber(accountNumber)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "not authorized"})
		return nil
	}

	// Get the 10 most recent transactions (transfers)
	transactions, err := s.store.GetTransactionHistory(acc.ID, 10, 0, "", 0)
	if err != nil {
		return err
	}

	var events []ActivityEvent

	// Add account created event
	events = append(events, ActivityEvent{
		Type:      "account_created",
		Title:     "Account Created",
		Details:   fmt.Sprintf("Your GoBank account was created and verified"),
		Status:    "completed",
		Timestamp: acc.CreatedAt,
	})

	// Process transactions into activity events
	for _, tx := range transactions {
		if tx.TransactionType == "transfer" {
			var eventType, title, details string
			var amount int64

			if tx.FromAccountID == acc.ID {
				// This user sent the transfer
				toAcc, err := s.store.GetAccountByID(tx.ToAccountID)
				if err != nil {
					log.Printf("could not fetch recipient account %d: %v", tx.ToAccountID, err)
					continue
				}
				recipientName := strings.TrimSpace(toAcc.FirstName + " " + toAcc.LastName)
				if recipientName == "" {
					recipientName = maskedAccountNumber(toAcc.Number)
				}
				eventType = "transfer_sent"
				title = fmt.Sprintf("Sent $%.2f to %s", float64(tx.Amount)/100.0, recipientName)
				details = fmt.Sprintf("Transfer ID: %s", tx.TransactionID)
				amount = tx.Amount
			} else {
				// This user received the transfer
				fromAcc, err := s.store.GetAccountByID(tx.FromAccountID)
				if err != nil {
					log.Printf("could not fetch sender account %d: %v", tx.FromAccountID, err)
					continue
				}
				senderName := strings.TrimSpace(fromAcc.FirstName + " " + fromAcc.LastName)
				if senderName == "" {
					senderName = maskedAccountNumber(fromAcc.Number)
				}
				eventType = "transfer_received"
				title = fmt.Sprintf("Received $%.2f from %s", float64(tx.Amount)/100.0, senderName)
				details = fmt.Sprintf("Transfer ID: %s", tx.TransactionID)
				amount = tx.Amount
			}

			events = append(events, ActivityEvent{
				Type:      eventType,
				Title:     title,
				Details:   details,
				Amount:    amount,
				Status:    tx.Status,
				Timestamp: tx.CreatedAt,
			})
		}
	}

	// Check if account has passkey and add event if it was recently created
	// (This is a simplified check; in a production system you'd track passkey registration times)
	if acc.HasPasskey {
		events = append(events, ActivityEvent{
			Type:      "passkey_registered",
			Title:     "Passkey Registered",
			Details:   "A passkey has been registered for secure login",
			Status:    "completed",
			Timestamp: acc.CreatedAt.Add(1 * time.Minute), // Placeholder - ideally we'd track this
		})
	}

	// Sort events by timestamp descending (newest first)
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.After(events[j].Timestamp)
	})

	// Return only the 5 most recent events for the dashboard preview
	if len(events) > 5 {
		events = events[:5]
	}

	return WriteJSON(w, http.StatusOK, events)
}

// parseMonth converts a month string (number 1-12 or English name) to its
// numeric representation (1-12). Returns 0 if the value cannot be parsed.
func parseMonth(m string) int {
	m = strings.ToLower(strings.TrimSpace(m))
	monthNames := map[string]int{
		"january": 1, "february": 2, "march": 3, "april": 4,
		"may": 5, "june": 6, "july": 7, "august": 8,
		"september": 9, "october": 10, "november": 11, "december": 12,
	}
	if n, ok := monthNames[m]; ok {
		return n
	}
	if n, err := strconv.Atoi(m); err == nil && n >= 1 && n <= 12 {
		return n
	}
	return 0
}

// getAuthenticatedAccount validates the JWT in the Authorization header and
// returns the corresponding account. It writes an appropriate HTTP error
// response and returns nil when authentication fails.
func (s *APIServer) getAuthenticatedAccount(w http.ResponseWriter, r *http.Request) *Account {
	tokenString := r.Header.Get("Authorization")
	token, err := validateJWT(tokenString, s.cfg.JWTSecret)
	if err != nil || !token.Valid {
		WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "not authorized"})
		return nil
	}

	claims := token.Claims.(jwt.MapClaims)
	accountNumber := int64(claims["accountNumber"].(float64))

	acc, err := s.store.GetAccountByNumber(accountNumber)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "not authorized"})
		return nil
	}
	return acc
}

// handleUpdateAccount handles POST /account/update.
// The request body must contain an "action" field that selects the update flow.
func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	acc := s.getAuthenticatedAccount(w, r)
	if acc == nil {
		return nil
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	switch req.Action {
	case "profile":
		return s.handleUpdateProfile(w, acc, &req)
	case "email_request":
		return s.handleEmailUpdateRequest(w, acc, &req)
	case "email_verify":
		return s.handleEmailUpdateVerify(w, acc, &req)
	case "password":
		return s.handlePasswordUpdate(w, acc, &req)
	default:
		return fmt.Errorf("unknown action: %s", req.Action)
	}
}

func (s *APIServer) handleUpdateProfile(w http.ResponseWriter, acc *Account, req *UpdateProfileRequest) error {
	if strings.TrimSpace(req.FirstName) == "" || strings.TrimSpace(req.LastName) == "" {
		return fmt.Errorf("firstName and lastName are required")
	}

	if err := s.store.UpdateAccountName(acc.ID, req.FirstName, req.LastName); err != nil {
		return err
	}

	acc.FirstName = req.FirstName
	acc.LastName = req.LastName

	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIServer) handleEmailUpdateRequest(w http.ResponseWriter, acc *Account, req *UpdateProfileRequest) error {
	if _, err := mail.ParseAddress(req.NewEmail); err != nil {
		return fmt.Errorf("a valid email address is required")
	}

	if !acc.ValidPassword(req.Password) {
		return fmt.Errorf("not authorized")
	}

	code, err := generateVerificationCode()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	ppu := &PendingProfileUpdate{
		AccountID:        acc.ID,
		NewEmail:         req.NewEmail,
		VerificationCode: code,
		ExpiresAt:        now.Add(5 * time.Minute),
		CreatedAt:        now,
	}

	if err := s.store.CreatePendingProfileUpdate(ppu); err != nil {
		return err
	}

	if err := sendVerificationEmail(req.NewEmail, code, s.cfg); err != nil {
		_ = s.store.DeletePendingProfileUpdate(ppu.ID)
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message": "verification code sent to new email address",
	})
}

func (s *APIServer) handleEmailUpdateVerify(w http.ResponseWriter, acc *Account, req *UpdateProfileRequest) error {
	if len(req.OTP) != 6 {
		return fmt.Errorf("OTP must be 6 digits")
	}

	if strings.TrimSpace(req.NewEmail) == "" {
		return fmt.Errorf("newEmail is required")
	}

	ppu, err := s.store.GetPendingProfileUpdate(acc.ID, req.NewEmail, req.OTP)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}

	if time.Now().UTC().After(ppu.ExpiresAt) {
		_ = s.store.DeletePendingProfileUpdate(ppu.ID)
		return fmt.Errorf("OTP has expired")
	}

	if err := s.store.UpdateAccountEmail(acc.ID, ppu.NewEmail); err != nil {
		return err
	}

	if err := s.store.DeletePendingProfileUpdate(ppu.ID); err != nil {
		log.Printf("failed to delete pending profile update %d: %v", ppu.ID, err)
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message": "email updated successfully",
	})
}

func (s *APIServer) handlePasswordUpdate(w http.ResponseWriter, acc *Account, req *UpdateProfileRequest) error {
	if !acc.ValidPassword(req.CurrentPassword) {
		return fmt.Errorf("not authorized")
	}

	if len(req.NewPassword) < 8 {
		return fmt.Errorf("new password must be at least 8 characters")
	}

	if req.NewPassword != req.ConfirmPassword {
		return fmt.Errorf("new password and confirm password do not match")
	}

	if acc.ValidPassword(req.NewPassword) {
		return fmt.Errorf("new password must differ from current password")
	}

	encpw, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.store.UpdateAccountPassword(acc.ID, string(encpw)); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message": "password updated successfully",
	})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	// Handle account transfer API requests here
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	tokenString := r.Header.Get("Authorization")
	token, err := validateJWT(tokenString, s.cfg.JWTSecret)
	if err != nil || !token.Valid {
		WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "not authorized"})
		return nil
	}

	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	if transferReq.Amount <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}

	claims := token.Claims.(jwt.MapClaims)
	fromAccountNumber := int64(claims["accountNumber"].(float64))

	fromAcc, err := s.store.GetAccountByNumber(fromAccountNumber)
	if err != nil {
		return fmt.Errorf("source account not found")
	}

	toAcc, err := s.store.GetAccountByNumber(transferReq.ToAccount)
	if err != nil {
		return fmt.Errorf("destination account not found")
	}

	if fromAcc.ID == toAcc.ID {
		return fmt.Errorf("cannot transfer to the same account")
	}

	transfer, err := s.store.Transfer(fromAcc.ID, toAcc.ID, transferReq.Amount)
	if err != nil {
		return err
	}

	senderName := strings.TrimSpace(fromAcc.FirstName + " " + fromAcc.LastName)
	if senderName == "" {
		senderName = maskedAccountNumber(fromAcc.Number)
	}

	recipientName := strings.TrimSpace(toAcc.FirstName + " " + toAcc.LastName)
	if recipientName == "" {
		recipientName = maskedAccountNumber(toAcc.Number)
	}

	// Send emails asynchronously
	go func(sender, recipient *Account, result *TransferResult, cfg *Config) {
		// Send sender (payer) notification
		if err := sendTransferNotificationEmail(sender.Email, false, senderName, recipientName, result, recipient, cfg); err != nil {
			log.Printf("failed to send sender receipt for transaction %s: %v", result.TransactionID, err)
		}

		// Send recipient (payee) notification
		if err := sendTransferNotificationEmail(recipient.Email, true, recipientName, senderName, result, recipient, cfg); err != nil {
			log.Printf("failed to send recipient receipt for transaction %s: %v", result.TransactionID, err)
		}
	}(fromAcc, toAcc, transfer, s.cfg)

	return WriteJSON(w, http.StatusOK, map[string]any{
		"status":        "transfer successful",
		"transactionId": transfer.TransactionID,
		"transferredAt": transfer.CreatedAt,
	})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func createJWT(account *Account, secret string) (string, error) {
	//create the claims
	claims := &jwt.MapClaims{
		"exp":           time.Now().Add(24 * time.Hour).Unix(),
		"accountNumber": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage, secret string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("we are calling the withJWTAuth middleware")

		tokenString := r.Header.Get("Authorization")

		token, err := validateJWT(tokenString, secret)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid token"})
			return
		}
		if !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid token"})
			return
		}

		userID, err := getID(r)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid token"})
			return
		}

		account, err := s.GetAccountByID(userID)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if int64(claims["accountNumber"].(float64)) != account.Number {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid token"})
			return
		}

		handlerFunc(w, r)
	}

}

func validateJWT(tokenString, secret string) (*jwt.Token, error) {
	tokenString = strings.TrimSpace(tokenString)
	if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
		tokenString = strings.TrimSpace(tokenString[7:])
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Handle the error, e.g., log it and return an appropriate HTTP response
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return id, fmt.Errorf("Invalid id given %s", idstr)
	}

	return id, nil
}
