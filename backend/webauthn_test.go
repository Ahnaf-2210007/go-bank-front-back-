package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockStore struct{}

func (m *mockStore) CreateAccount(*Account) error {
	return nil
}
func (m *mockStore) GetAccountByID(int) (*Account, error) {
	return nil, nil
}
func (m *mockStore) DeleteAccount(int) error {
	return nil
}
func (m *mockStore) GetAccounts() ([]*Account, error) {
	return nil, nil
}
func (m *mockStore) UpdateAccount(*Account) error {
	return nil
}
func (m *mockStore) GetAccountByNumber(number int64) (*Account, error) {
	return &Account{ID: 1, Number: number, Email: "test@test.com"}, nil
}
func (m *mockStore) Transfer(fromID, toID int, amount int64) (*TransferResult, error) {
	return nil, nil
}
func (m *mockStore) RedeemCoupon(accountID int, couponCode string) error {
	return nil
}
func (m *mockStore) CreatePendingAccount(acc *PendingAccount) error {
	return nil
}
func (m *mockStore) GetPendingAccountByCode(code string) (*PendingAccount, error) {
	return nil, nil
}
func (m *mockStore) DeletePendingAccount(id int) error {
	return nil
}
func (m *mockStore) GetTransactionHistory(accountID, limit, offset int, txType string, month int) ([]*TransactionRecord, error) {
	return nil, nil
}
func (m *mockStore) CreatePendingProfileUpdate(ppu *PendingProfileUpdate) error {
	return nil
}
func (m *mockStore) GetPendingProfileUpdate(accountID int, newEmail, code string) (*PendingProfileUpdate, error) {
	return nil, nil
}
func (m *mockStore) DeletePendingProfileUpdate(id int) error {
	return nil
}
func (m *mockStore) UpdateAccountName(accountID int, firstName, lastName string) error {
	return nil
}
func (m *mockStore) UpdateAccountEmail(accountID int, newEmail string) error {
	return nil
}
func (m *mockStore) UpdateAccountPassword(accountID int, encryptedPassword string) error {
	return nil
}
func (m *mockStore) GetAccountByEmail(email string) (*Account, error) {
	return &Account{ID: 1, Email: email}, nil
}
func (m *mockStore) CreateWebAuthnCredential(cred *WebAuthnCredential, accountID int) error {
	return nil
}
func (m *mockStore) GetWebAuthnCredentialsByAccountID(accountID int) ([]*WebAuthnCredential, error) {
	return nil, nil
}
func (m *mockStore) CreateWebAuthnCredentialTable() error {
	return nil
}
func (m *mockStore) UpdateAccountPasskeyStatus(accountID int, hasPasskey bool) error {
	return nil
}

func TestWebAuthnHandlers(t *testing.T) {
	store := &mockStore{}
	cfg := &Config{JWTSecret: "test-secret"}
	server := NewAPIServer(":3000", store, cfg)

	t.Run("handleRegisterBegin", func(t *testing.T) {
		acc, _ := NewAccount("test", "test", "test@test.com", "password")
		token, _ := createJWT(acc, cfg.JWTSecret)
		req, _ := http.NewRequest("POST", "/webauthn/register/begin", nil)
		req.Header.Set("Authorization", token)
		rr := httptest.NewRecorder()
		handler := makeHTTPHandlerFunc(server.webAuthnHandler.handleRegisterBegin)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("handleLoginBegin", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/webauthn/login/begin", strings.NewReader(`{"email": "test@test.com"}`))
		rr := httptest.NewRecorder()
		handler := makeHTTPHandlerFunc(server.webAuthnHandler.handleLoginBegin)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
