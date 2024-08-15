package Infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "example/go-clean-architecture/Domain"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestAuthMiddleware_ValidToken tests the AuthMiddleware function when a valid token is provided.
// It creates a valid user instance and generates a valid token for the user.
// Then, it creates a request with the valid token and sends it to the "/test" endpoint.
// Finally, it asserts that the response code is 200 and the response body contains the expected JSON message.
func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Create a valid user instance
	user := domain.User{
		ID:       primitive.NewObjectID(), // Replace with actual MongoDB ObjectID.Hex() if applicable
		Username: "testuser",
		Role:     "USER",
	}

	// Generate a valid token for the user
	tokenString, err := GenerateToken(user)
	assert.Nil(t, err)

	// Create a request with the valid token
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer " + tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message":"Success"}`, w.Body.String())
}

// TestAuthMiddleware_InvalidToken tests the behavior of the AuthMiddleware function when an invalid token is provided.
// It creates a request with an invalid token and sends it to the server. The server should respond with a
// 401 Unauthorized status code and an error message indicating that the token is invalid.
// The test asserts that the response code is 401 and the response body is `{"error":"Invalid token"}`.
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Create a request with an invalid token
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.JSONEq(t, `{"error":"Invalid token"}`, w.Body.String())
}

// TestAuthMiddleware_ExpiredToken tests the behavior of the AuthMiddleware function when an expired token is provided.
//
// It creates an expired token and sets it as the Authorization header in a request. 
//Then, it sends the request to the "/test" endpoint and checks if the response code is 401 Unauthorized and the response body contains the error message "Token has expired".
// This test ensures that the AuthMiddleware function correctly handles expired tokens and returns the appropriate error response.
func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// Create an expired token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": primitive.NewObjectID().Hex(),
		"role":    "USER",
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(SECRET_KEY)

	// Create a request with the expired token
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.JSONEq(t, `{"error":"Token has expired"}`, w.Body.String())
}

// TestAuthAdminMiddleware_ValidAdmin tests the AuthAdminMiddleware function when a valid "ADMIN" role is set in the context.
//
// It creates a request with an "ADMIN" role in the context and sends a GET request to "/admin" endpoint.
// The expected response is a JSON object with a message indicating that admin access is granted.
//
// This test asserts that the HTTP status code is 200 and the response body matches the expected JSON.
func TestAuthAdminMiddleware_ValidAdmin(t *testing.T) {
	// Create a request with an "ADMIN" role in the context
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.Use(func(c *gin.Context) {
		c.Set("role", "ADMIN")
		c.Next()
	})
	r.Use(AuthAdminMiddleware())
	r.GET("/admin", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Admin access granted"})
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message":"Admin access granted"}`, w.Body.String())
}

// TestAuthAdminMiddleware_InvalidRole tests the behavior of the AuthAdminMiddleware function when the user has an invalid role.
//
// This test creates a request with a non-admin role in the context and sends it to the "/admin" endpoint.
// The AuthAdminMiddleware function is used as a middleware to check if the user has admin access.
// The expected behavior is that the response status code should be 403 (Forbidden) and the response body should contain the error message "Admin access required".
func TestAuthAdminMiddleware_InvalidRole(t *testing.T) {
	// Create a request with a non-admin role in the context
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.Use(func(c *gin.Context) {
		c.Set("role", "USER")
		c.Next()
	})
	r.Use(AuthAdminMiddleware())
	r.GET("/admin", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Admin access granted"})
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
	assert.JSONEq(t, `{"error":"Admin access required"}`, w.Body.String())
}
