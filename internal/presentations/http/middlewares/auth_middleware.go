package middlewares

import (
	"backend/internal/constants/enums/jwtpurpose"
	"backend/pkg/utils/jwt"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthHeader(secret []byte, purpose jwtpurpose.JWTPurpose) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("missing Authorization header")
			permissionDenied(c)
			return
		}

		// split bearer
		if !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("invalid Authorization header format")
			permissionDenied(c)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// validate token
		claims, err := jwt.ValidateToken(secret, token, purpose)
		if err != nil {
			fmt.Println("failed to validate token", err)
			permissionDenied(c)
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			fmt.Println("error when parsing claims subject to uuid")
			permissionDenied(c)
			return
		}
		c.Set("userID", userID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func AuthQuery(secret []byte, purpose jwtpurpose.JWTPurpose) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			fmt.Println("missing token in query")
			permissionDenied(c)
			return
		}

		// validate token
		claims, err := jwt.ValidateToken(secret, token, purpose)
		if err != nil {
			fmt.Println("failed to validate token", err)
			permissionDenied(c)
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			fmt.Println("error when parsing claims subject to uuid")
			permissionDenied(c)
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func permissionDenied(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "permission denied",
	})
}

// RequireRolesFromToken checks allowed roles from token claims (no DB)
func RequireRolesFromToken(allowedRoles ...string) gin.HandlerFunc {
	roleSet := map[string]struct{}{}
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}
	return func(c *gin.Context) {
		val, ok := c.Get("role")
		if !ok {
			permissionDenied(c)
			return
		}
		roleStr, ok := val.(string)
		if !ok {
			permissionDenied(c)
			return
		}
		if _, ok := roleSet[roleStr]; !ok {
			permissionDenied(c)
			return
		}
		c.Next()
	}
}
