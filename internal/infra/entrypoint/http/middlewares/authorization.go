package middlewares

import (
	"net/http"
	"slices"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	sessionManager ports.SessionManager
}

func NewAuthMiddleware(sessionManager ports.SessionManager) *AuthMiddleware {
	return &AuthMiddleware{
		sessionManager: sessionManager,
	}
}

var NO_AUTH_NEEDED = []string{
	"/login",
	"/register",
}

func (middleware *AuthMiddleware) isPublicRoute(route string) bool {
	contain := slices.Contains(NO_AUTH_NEEDED, route)
	return contain
}

func (middleware *AuthMiddleware) CheckAuthMiddleware(c *gin.Context) {
	path := c.Request.URL.Path
	if middleware.isPublicRoute(path) {
		c.Next()
		return
	}
	authCookie, err := c.Request.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Cannot validate the Authorization")
		c.Abort()
		return
	}
	_, err = middleware.sessionManager.GetCredentials(authCookie.Value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"err": "Cannot validate the Authorization",
			"c":   err.Error(),
		})
		c.Abort()
		return
	}
	c.Next()
}
