package middleware

import (
	"fmt"
	"hardware_store/internal/web/handler/auth"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type contextKey string

const uidKey contextKey = "user_id"
const appKey contextKey = "app_id"

func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	fmt.Printf("DEBUG RAW HEADER: [%s]\n", authHeader)
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) < 2 {
		return ""
	}
	return splitToken[1]
}

func NewMiddleware(log *slog.Logger, auth auth.TokenService) gin.HandlerFunc {
	const op = "middleware.NewMiddleware"
	log = log.With(slog.String("op", op))

	return func(c *gin.Context) {
		token := extractBearerToken(c)
		fmt.Printf("DEBUG CLEAN TOKEN: [%s]\n", token)
		if token == "" {
			log.Info("No token provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		userID, appID, err := auth.ValidateToken(c.Request.Context(), token)
		if err != nil {
			log.Info("Invalid token", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		log.Debug("User authorized",
			slog.String("user_id", userID),
			slog.String("app_id", appID),
		)

		c.Set(string(uidKey), userID)
		c.Set(string(appKey), appID)

		c.Next()
	}
}
