package Infrastructure

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// checks the users password with the one in the database
func VerifyPassword(userPassword string, foundPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(foundPassword), []byte(userPassword))
	return err == nil
}