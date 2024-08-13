package Infrastructure

import (
	domain "example/go-clean-architecture/Domain"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SECRET_KEY = []byte("MY-Secret-Key")

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateTokens generates a JWT token for the given user.
// It takes a user object as input and returns the generated token string and an error, if any.
// The token is valid for 24 hours from the current time.
// The user's ID, username, and role are included as claims in the token.
// The token is signed using the SECRET_KEY.
func GenerateToken(user domain.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	claims := &Claims{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	//generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
