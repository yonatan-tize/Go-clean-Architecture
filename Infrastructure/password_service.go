package Infrastructure

import "golang.org/x/crypto/bcrypt"


// HashPassword takes a password string and returns the hashed version of the password.
// It uses bcrypt.GenerateFromPassword to generate the hash with the default cost.
// If an error occurs during the hashing process, it returns an empty string and the error.
// Otherwise, it returns the hashed password as a string and nil error.
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