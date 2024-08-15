package Infrastructure

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware function that handles authentication for incoming requests.
// It checks the "Authorization" header in the request and validates the token.
// If the header is missing or the token is invalid, it returns a 401 Unauthorized response.
// If the token has expired, it returns a 401 Unauthorized response.
// If the token is valid, it sets the "user_id" and "role" values in the context and allows the request to proceed.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(authParts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte("MY-Secret-Key"), nil
		})

		if err != nil {
			// Check if the error is due to token expiration
			if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors == jwt.ValidationErrorExpired {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			c.Set("role", claims["role"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}


// AuthAdminMiddleware is a middleware function that checks if the user has admin access.
// It retrieves the user's role from the context and verifies if it is set to "ADMIN".
// If the role is not set or is not "ADMIN", it returns a JSON response with a 403 Forbidden status
// and an error message indicating that admin access is required.
// If the role is "ADMIN", it allows the request to proceed to the next middleware or handler.
func AuthAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
