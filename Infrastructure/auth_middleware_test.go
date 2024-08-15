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
