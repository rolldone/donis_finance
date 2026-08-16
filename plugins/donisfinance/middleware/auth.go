package middleware

import (
	"net/http"
	"strings"

	"github.com/rolldone/donisgo/internal/auth"

	"github.com/gin-gonic/gin"
)

// JWTAuth returns a middleware that verifies JWT access tokens
// and injects user_id + role into the gin context.
// The role parameter is the expected role ("admin" or "member").
func JWTAuth(expectedRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format, expected 'Bearer <token>'"})
			return
		}

		tokenStr := parts[1]
		userID, err := auth.ParseAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Inject into context for downstream handlers
		c.Set("user_id", userID)
		c.Set("role", expectedRole)

		c.Next()
	}
}

// ExtractAuth is a lightweight middleware that optionally parses the token
// if present, but doesn't block unauthenticated requests.
// Useful for endpoints that behave differently for auth'd users.
func ExtractAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.Next()
			return
		}

		userID, err := auth.ParseAccessToken(parts[1])
		if err == nil && userID != "" {
			c.Set("user_id", userID)
		}

		c.Next()
	}
}
