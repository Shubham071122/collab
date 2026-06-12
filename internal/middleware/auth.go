package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/shubham071122/collab/internal/config"
	"github.com/shubham071122/collab/internal/response"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		cfg := config.LoadConfig()
		jwtSecret := []byte(cfg.JWTSecret)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			tokenQuery := c.Query("token")
			if tokenQuery != "" {
				authHeader = "Bearer " + tokenQuery
			}
		}
		if authHeader == "" {
			if cookieToken, err := c.Cookie("auth_token"); err == nil && cookieToken != "" {
				authHeader = "Bearer " + cookieToken
			}
		}
		if authHeader == "" {
			protocol := c.GetHeader("Sec-WebSocket-Protocol")
			if protocol != "" {
				parts := strings.Split(protocol, ",")
				for _, part := range parts {
					trimmed := strings.TrimSpace(part)
					if trimmed != "" && trimmed != "Bearer" {
						authHeader = "Bearer " + trimmed
						c.Header("Sec-WebSocket-Protocol", trimmed)
						break
					}
				}
			}
		}

		if authHeader == "" {
			response.JSON(c, response.StatusUnauthorized, "Authorization header missing", nil, nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			},
		)

		if err != nil || !token.Valid {
			response.JSON(c, response.StatusUnauthorized, "Invalid token", nil, nil)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.JSON(c, response.StatusUnauthorized, "Invalid token claims", nil, nil)
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			response.JSON(c, response.StatusUnauthorized, "Invalid token claims", nil, nil)
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
