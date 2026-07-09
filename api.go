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

	msg := strings.Join([]string{
		"From: " + cfg.SMTPEmail,
		"To: " + email,
		"Subject: Your GoBank Verification Code",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"UTF-8\"",
		"",
		fmt.Sprintf("Your GoBank verification code is: %s", code),
		"This code expires in 5 minutes.",
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

func sendTransferNotificationEmail(email, subject string, lines []string, transfer *TransferResult, cfg *Config) error {
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

	host := cfg.SMTPHost
	port := cfg.SMTPPort
	addr := host + ":" + port
	auth := smtp.PlainAuth("", cfg.SMTPEmail, cfg.SMTPPassword, host)

	msg := strings.Join([]string{
		"From: " + cfg.SMTPEmail,
		"To: " + email,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"UTF-8\"",
		"",
	}, "\r\n")
	msg += "\r\n" + strings.Join(lines, "\r\n")

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

	msg := strings.Join([]string{
		"From: " + cfg.SMTPEmail,
		"To: " + account.Email,
		"Subject: Welcome to GoBank - Account Confirmation",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"UTF-8\"",
		"",
		fmt.Sprintf("Hello %s,", fullName),
		"",
		"Welcome to GoBank! Your account has been successfully created and verified.",
		"",
		"=== ACCOUNT DETAILS ===",
		fmt.Sprintf("Account Holder: %s", fullName),
		fmt.Sprintf("Account Number: %d", account.Number),
		fmt.Sprintf("Account ID: %d", account.ID),
		fmt.Sprintf("Current Balance: $%.2f", float64(account.Balance)/100),
		"",
		"=== NEXT STEPS ===",
		"1. Log in to your GoBank account",
		"2. Complete your profile with additional information",
		"3. Set up a passkey for enhanced security",
		"4. Start making transfers and enjoying exclusive offers",
		"",
		"=== SECURITY REMINDER ===",
		"- Never share your password or verification codes",
		"- Enable two-factor authentication for added security",
		"- Always use secure connections when accessing your account",
		"",
		"If you did not create this account, please contact our support team immediately.",
		"",
		"Best regards,",
		"The GoBank Team",
		"",
		"This is an automated message. Please do not reply to this email.",
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

	senderLines := []string{
		fmt.Sprintf("Hello %s,", senderName),
		"",
		"Your transfer has been completed successfully.",
		fmt.Sprintf("Transaction ID: %s", transfer.TransactionID),
		fmt.Sprintf("Recipient: %s", recipientName),
		fmt.Sprintf("Amount: %d", transfer.Amount),
		fmt.Sprintf("Date: %s", transfer.CreatedAt.Format(time.RFC1123Z)),
	}

	recipientLines := []string{
		fmt.Sprintf("Hello %s,", recipientName),
		"",
		"You have received a transfer successfully.",
		fmt.Sprintf("Transaction ID: %s", transfer.TransactionID),
		fmt.Sprintf("Sender: %s", senderName),
		fmt.Sprintf("Amount: %d", transfer.Amount),
		fmt.Sprintf("Date: %s", transfer.CreatedAt.Format(time.RFC1123Z)),
	}

	go func(sender, recipient *Account, result *TransferResult, senderBody, recipientBody []string, cfg *Config) {
		if err := sendTransferNotificationEmail(sender.Email, "GoBank Transfer Receipt", senderBody, result, cfg); err != nil {
			log.Printf("failed to send sender receipt for transaction %s: %v", result.TransactionID, err)
		}

		if err := sendTransferNotificationEmail(recipient.Email, "GoBank Incoming Transfer", recipientBody, result, cfg); err != nil {
			log.Printf("failed to send recipient receipt for transaction %s: %v", result.TransactionID, err)
		}
	}(fromAcc, toAcc, transfer, senderLines, recipientLines, s.cfg)

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
