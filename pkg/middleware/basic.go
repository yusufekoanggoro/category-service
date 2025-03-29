package middleware

import (
	"category-service/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware(cfg config.ConfigProvider) gin.HandlerFunc {
	username := cfg.GetBasicAuthUsername()
	password := cfg.GetBasicAuthPassword()

	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != username || pass != password {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
