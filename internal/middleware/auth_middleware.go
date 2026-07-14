package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/carloscfgos1980/graphql-habit-tracker/internal/utils"

	"github.com/gin-gonic/gin"
)

type contextKey string

// contextKey("userID") != string("userID")
const userIDKey contextKey = "userID"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization: "Bearer ..."
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := utils.ValidateJWT(tokenString, os.Getenv("JWT_SECRET"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), userIDKey, userID)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetUserID(ctx context.Context) (string, bool) {
	// interface{} == any{}
	userID, ok := ctx.Value(userIDKey).(string)

	return userID, ok
}
