package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"golang.org/x/crypto/bcrypt"
)

type WebAuthnCredential struct {
	ID              []byte
	PublicKey       []byte
	AttestationType string
	Transport       []string
	Flags           webauthn.CredentialFlags
	Authenticator   webauthn.Authenticator
}

func (wc *WebAuthnCredential) WebAuthnCredentialID() []byte {
	return wc.ID
}

func (wc *WebAuthnCredential) WebAuthnCredentialPublicKey() []byte {
	return wc.PublicKey
}

func (wc *WebAuthnCredential) WebAuthnCredentialAttestationType() string {
	return wc.AttestationType
}

func (wc *WebAuthnCredential) WebAuthnCredentialTransport() []protocol.AuthenticatorTransport {
	transports := make([]protocol.AuthenticatorTransport, len(wc.Transport))
	for i, t := range wc.Transport {
		transports[i] = protocol.AuthenticatorTransport(t)
	}
	return transports
}

func (wc *WebAuthnCredential) WebAuthnCredentialFlags() webauthn.CredentialFlags {
	return wc.Flags
}

func (wc *WebAuthnCredential) WebAuthnAuthenticator() webauthn.Authenticator {
	return wc.Authenticator
}

type LoginResponse struct {
	Number int64  `json:"number"`
	Token  string `json:"token"`
}

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type TransferRequest struct {
	ToAccount int64 `json:"toAccount"`
	Amount    int64 `json:"amount"`
}

type TransferResult struct {
	ID            int       `json:"id"`
	TransactionID string    `json:"transactionId"`
	FromAccountID int       `json:"fromAccountId"`
	ToAccountID   int       `json:"toAccountId"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"createdAt"`
}

// TransactionRecord represents a single entry in a user's transaction history.
type TransactionRecord struct {
	ID              int       `json:"id"`
	TransactionID   string    `json:"transactionId"`
	FromAccountID   int       `json:"fromAccountId"`
	ToAccountID     int       `json:"toAccountId"`
	Amount          int64     `json:"amount"`
	TransactionType string    `json:"transactionType"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

type VerificationRequest struct {
	Code string `json:"code"`
}

type OfferRequest struct {
	CouponCode string `json:"couponCode"`
}

// UpdateProfileRequest is the request body for POST /account/update.
// The Action field selects the update flow; the remaining fields are
// required only for the matching action.
type UpdateProfileRequest struct {
	Action string `json:"action"`

	// profile action
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	// email_request action
	NewEmail string `json:"newEmail"`
	Password string `json:"password"`

	// email_verify action
	OTP string `json:"otp"`

	// password action
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

// PendingProfileUpdate holds a pending email-change request awaiting OTP
// confirmation.  Rows expire after 5 minutes.
type PendingProfileUpdate struct {
	ID               int
	AccountID        int
	NewEmail         string
	VerificationCode string
	ExpiresAt        time.Time
	CreatedAt        time.Time
}

type Account struct {
	ID                  int                   `json:"id"`
	FirstName           string                `json:"firstName"`
	LastName            string                `json:"lastName"`
	Number              int64                 `json:"number"`
	Email               string                `json:"email"`
	EncryptedPassword   string                `json:"-"`
	Balance             int64                 `json:"balance"`
	HasPasskey          bool                  `json:"hasPasskey"`
	CreatedAt           time.Time             `json:"createdAt"`
	webAuthnCredentials []webauthn.Credential `json:"-"`
}

func (a *Account) WebAuthnID() []byte {
	return []byte(fmt.Sprintf("%d", a.ID))
}

func (a *Account) WebAuthnName() string {
	return a.Email
}

func (a *Account) WebAuthnDisplayName() string {
	return fmt.Sprintf("%s %s", a.FirstName, a.LastName)
}

func (a *Account) WebAuthnIcon() string {
	return ""
}

func (a *Account) WebAuthnCredentials() []webauthn.Credential {
	return a.webAuthnCredentials
}

// PendingAccount holds a newly registered account awaiting email verification.
type PendingAccount struct {
	ID                int
	FirstName         string
	LastName          string
	Email             string
	Number            int64
	EncryptedPassword string
	VerificationCode  string
	ExpiresAt         time.Time
}

func (a *Account) ValidPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password)) == nil
}

func NewAccount(firstName, lastName, email, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Email:             email,
		Number:            n.Int64(),
		EncryptedPassword: string(encpw),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
