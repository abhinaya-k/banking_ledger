package middleware

import (
	"banking_ledger/config"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Your JWT secret key (should come from env ideally)
var jwtSecret = []byte(config.JWT_SECRET) // Replace with your env/config

func AuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			c.Abort()
			return
		}

		// Check format: "Bearer <token>"
		var tokenString string
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 || strings.ToLower(tokenParts[0]) == "bearer" {
			tokenString = tokenParts[1]
		} else {
			tokenString = tokenParts[0]
		}

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Extract user ID from the token
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			sub, ok := claims["sub"].(string) // jwt parses numbers as float64
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
				c.Abort()
				return
			}
			userID, err := strconv.Atoi(sub)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sub typ"})
				c.Abort()
				return
			}
			// Set userID in the context
			c.Set("user_id", userID)

			role, ok := claims["role"].(string) // jwt parses numbers as float64
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims for role"})
				c.Abort()
				return
			}
			name, ok := claims["role"].(string) // jwt parses numbers as float64
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims for name"})
				c.Abort()
				return
			}

			c.Set("role", role)
			c.Set("name", name)

			// Continue to next handler
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}
	}
}
