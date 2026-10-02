package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const publicRSAKey = "../publicKey.pem"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization required"})
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		pkBytes, err := os.ReadFile(publicRSAKey)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error reading public key"})
			return
		}

		pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pkBytes)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error parsing public key"})
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return pubKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token or expired"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
