package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("John", "Doe", "john@example.com", "password123")
	assert.Nil(t, err)

	fmt.Printf("%+v\n", acc)
}

func TestNewAccountPasswordIsHashed(t *testing.T) {
	acc, err := NewAccount("Jane", "Smith", "jane@example.com", "securepass")
	assert.Nil(t, err)
	assert.True(t, acc.ValidPassword("securepass"))
	assert.False(t, acc.ValidPassword("wrongpassword"))
}

func TestNewAccountNumberInRange(t *testing.T) {
	for i := 0; i < 10; i++ {
		acc, err := NewAccount("Test", "User", "test@example.com", "testpassword")
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, acc.Number, int64(0))
		assert.Less(t, acc.Number, int64(1000000))
	}
}

func TestUpdateAccountPasswordHashing(t *testing.T) {
	newPassword := "newSecurePass99"
	encpw, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	assert.Nil(t, err)
	assert.NotEmpty(t, encpw)

	acc := &Account{EncryptedPassword: string(encpw)}
	assert.True(t, acc.ValidPassword(newPassword))
	assert.False(t, acc.ValidPassword("wrongPassword"))
}

func TestUpdateProfileRequestFields(t *testing.T) {
	req := UpdateProfileRequest{
		Action:          "password",
		CurrentPassword: "oldPass123",
		NewPassword:     "newPass456",
		ConfirmPassword: "newPass456",
	}
	assert.Equal(t, "password", req.Action)
	assert.Equal(t, "oldPass123", req.CurrentPassword)
	assert.Equal(t, "newPass456", req.NewPassword)
	assert.Equal(t, req.NewPassword, req.ConfirmPassword)
}

func TestPendingProfileUpdateStruct(t *testing.T) {
	ppu := PendingProfileUpdate{
		AccountID:        1,
		NewEmail:         "new@example.com",
		VerificationCode: "123456",
	}
	assert.Equal(t, 1, ppu.AccountID)
	assert.Equal(t, "new@example.com", ppu.NewEmail)
	assert.Equal(t, "123456", ppu.VerificationCode)
}
