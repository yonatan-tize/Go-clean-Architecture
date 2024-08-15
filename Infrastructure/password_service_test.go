package Infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// TestHashPassword_Success tests that a password is hashed successfully
func TestHashPassword_Success(t *testing.T) {
	password := "securepassword"

	hashedPassword, err := HashPassword(password)
	assert.Nil(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Ensure that the hashed password is different from the plain password
	assert.NotEqual(t, password, hashedPassword)

	// Verify that the hashed password matches the original password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.Nil(t, err)
}


// TestVerifyPassword_Success tests that the password verification succeeds for matching passwords
func TestVerifyPassword_Success(t *testing.T) {
	password := "securepassword"

	hashedPassword, err := HashPassword(password)
	assert.Nil(t, err)

	// Verify that the original password matches the hashed password
	isValid := VerifyPassword(password, hashedPassword)
	assert.True(t, isValid)
}

// TestVerifyPassword_Failure tests that the password verification fails for non-matching passwords
func TestVerifyPassword_Failure(t *testing.T) {
	password := "securepassword"
	wrongPassword := "wrongpassword"

	hashedPassword, err := HashPassword(password)
	assert.Nil(t, err)

	// Verify that the wrong password does not match the hashed password
	isValid := VerifyPassword(wrongPassword, hashedPassword)
	assert.False(t, isValid)
}
