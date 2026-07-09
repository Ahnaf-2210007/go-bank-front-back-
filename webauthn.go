package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// For this example, we'll use a simple in-memory store for session data.
// In a real-world application, you would use a more robust session management system.
var sessionStore = make(map[string]*webauthn.SessionData)

type WebAuthnHandler struct {
	store    Storage
	webauthn *webauthn.WebAuthn
	cfg      *Config
}

func NewWebAuthnHandler(store Storage, cfg *Config) (*WebAuthnHandler, error) {
	w, err := webauthn.New(&webauthn.Config{
		RPDisplayName: cfg.WebAuthnDisplayName,
		RPID:          cfg.WebAuthnRPID,
		RPOrigins:     []string{cfg.WebAuthnRPOrigin},
	})
	if err != nil {
		return nil, err
	}

	return &WebAuthnHandler{
		store:    store,
		webauthn: w,
		cfg:      cfg,
	}, nil
}

// getAuthenticatedAccount validates the JWT in the Authorization header and
// returns the corresponding account. It returns an error if authentication fails.
func (h *WebAuthnHandler) getAuthenticatedAccount(r *http.Request) (*Account, error) {
	tokenString := r.Header.Get("Authorization")
	token, err := validateJWT(tokenString, h.cfg.JWTSecret)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("not authorized")
	}

	claims := token.Claims.(jwt.MapClaims)
	accountNumber := int64(claims["accountNumber"].(float64))

	account, err := h.store.GetAccountByNumber(accountNumber)
	if err != nil {
		return nil, fmt.Errorf("not authorized")
	}
	return account, nil
}

func (h *WebAuthnHandler) handleRegisterBegin(w http.ResponseWriter, r *http.Request) error {
	account, err := h.getAuthenticatedAccount(r)
	if err != nil {
		return err
	}

	// Get the user's existing credentials from the database
	webAuthnCredentials, err := h.store.GetWebAuthnCredentialsByAccountID(account.ID)
	if err != nil {
		return err
	}
	var creds []webauthn.Credential
	for _, c := range webAuthnCredentials {
		// 1. Map the []string transports back to []protocol.AuthenticatorTransport
		transports := make([]protocol.AuthenticatorTransport, len(c.Transport))
		for i, t := range c.Transport {
			transports[i] = protocol.AuthenticatorTransport(t)
		}

		// 2. Map your custom DB struct back to the library's struct
		libCred := webauthn.Credential{
			ID:              c.ID,
			PublicKey:       c.PublicKey,
			AttestationType: c.AttestationType,
			Transport:       transports,
			Flags:           c.Flags,
			Authenticator:   c.Authenticator,
		}

		creds = append(creds, libCred)
	}
	account.webAuthnCredentials = creds

	options, sessionData, err := h.webauthn.BeginRegistration(
		account,
	)
	if err != nil {
		return err
	}

	// Store the session data for the registration process
	sessionStore[account.Email] = sessionData

	return WriteJSON(w, http.StatusOK, options)
}

func (h *WebAuthnHandler) handleRegisterFinish(w http.ResponseWriter, r *http.Request) error {
	account, err := h.getAuthenticatedAccount(r)
	if err != nil {
		return err
	}

	// Get the user's existing credentials from the database
	webAuthnCredentials, err := h.store.GetWebAuthnCredentialsByAccountID(account.ID)
	if err != nil {
		return err
	}
	var creds []webauthn.Credential
	for _, c := range webAuthnCredentials {
		// 1. Map the []string transports back to []protocol.AuthenticatorTransport
		transports := make([]protocol.AuthenticatorTransport, len(c.Transport))
		for i, t := range c.Transport {
			transports[i] = protocol.AuthenticatorTransport(t)
		}

		// 2. Map your custom DB struct back to the library's struct
		libCred := webauthn.Credential{
			ID:              c.ID,
			PublicKey:       c.PublicKey,
			AttestationType: c.AttestationType,
			Transport:       transports,
			Flags:           c.Flags,
			Authenticator:   c.Authenticator,
		}

		creds = append(creds, libCred)
	}
	account.webAuthnCredentials = creds

	// Get the session data that was stored during registration
	sessionData, ok := sessionStore[account.Email]
	if !ok {
		return fmt.Errorf("no session data for user")
	}

	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(r.Body)
	if err != nil {
		return err
	}

	credential, err := h.webauthn.CreateCredential(account, *sessionData, parsedResponse)
	if err != nil {
		return err
	}

	transports := make([]string, len(credential.Transport))
	for i, t := range credential.Transport {
		transports[i] = string(t)
	}

	webAuthnCred := &WebAuthnCredential{
		ID:              credential.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		Transport:       transports,
		Flags:           credential.Flags,
		Authenticator:   credential.Authenticator,
	}

	if err := h.store.CreateWebAuthnCredential(webAuthnCred, account.ID); err != nil {
		return err
	}

	// Update account to mark hasPasskey as true
	account.HasPasskey = true
	if err := h.store.UpdateAccountPasskeyStatus(account.ID, true); err != nil {
		log.Printf("failed to update passkey status for account %d: %v", account.ID, err)
		// Don't fail the registration response, just log the error
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *WebAuthnHandler) handleLoginBegin(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	account, err := h.store.GetAccountByEmail(req.Email)
	if err != nil {
		return err
	}

	webAuthnCredentials, err := h.store.GetWebAuthnCredentialsByAccountID(account.ID)
	if err != nil {
		return err
	}
	var creds []webauthn.Credential
	for _, c := range webAuthnCredentials {
		// 1. Map the []string transports back to []protocol.AuthenticatorTransport
		transports := make([]protocol.AuthenticatorTransport, len(c.Transport))
		for i, t := range c.Transport {
			transports[i] = protocol.AuthenticatorTransport(t)
		}

		// 2. Map your custom DB struct back to the library's struct
		libCred := webauthn.Credential{
			ID:              c.ID,
			PublicKey:       c.PublicKey,
			AttestationType: c.AttestationType,
			Transport:       transports,
			Flags:           c.Flags,
			Authenticator:   c.Authenticator,
		}

		creds = append(creds, libCred)
	}
	account.webAuthnCredentials = creds

	options, sessionData, err := h.webauthn.BeginLogin(account)
	if err != nil {
		return err
	}

	// Store the session data for the login process
	sessionStore[account.Email] = sessionData

	return WriteJSON(w, http.StatusOK, options)
}

func (h *WebAuthnHandler) handleLoginFinish(w http.ResponseWriter, r *http.Request) error {
	email, ok := mux.Vars(r)["email"]
	if !ok {
		return fmt.Errorf("email not provided in URL")
	}

	account, err := h.store.GetAccountByEmail(email)
	if err != nil {
		return err
	}

	// Get the user's existing credentials from the database
	webAuthnCredentials, err := h.store.GetWebAuthnCredentialsByAccountID(account.ID)
	if err != nil {
		return err
	}
	var creds []webauthn.Credential
	for _, c := range webAuthnCredentials {
		// 1. Map the []string transports back to []protocol.AuthenticatorTransport
		transports := make([]protocol.AuthenticatorTransport, len(c.Transport))
		for i, t := range c.Transport {
			transports[i] = protocol.AuthenticatorTransport(t)
		}

		// 2. Map your custom DB struct back to the library's struct
		libCred := webauthn.Credential{
			ID:              c.ID,
			PublicKey:       c.PublicKey,
			AttestationType: c.AttestationType,
			Transport:       transports,
			Flags:           c.Flags,
			Authenticator:   c.Authenticator,
		}

		creds = append(creds, libCred)
	}
	account.webAuthnCredentials = creds

	// Get the session data that was stored during login
	sessionData, ok := sessionStore[account.Email]
	if !ok {
		return fmt.Errorf("no session data for user")
	}

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(r.Body)
	if err != nil {
		return err
	}

	_, err = h.webauthn.ValidateLogin(account, *sessionData, parsedResponse)
	if err != nil {
		return err
	}

	token, err := createJWT(account, h.cfg.JWTSecret)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Token:  token,
		Number: account.Number,
	}

	return WriteJSON(w, http.StatusOK, resp)
}
