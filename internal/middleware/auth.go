package middleware

import (
	"github.com/gin-gonic/gin"
)

func Auth(r *gin.Engine, basicAuth *string) {
	r.Use(func(c *gin.Context) {
		if c.Request.Method == "GET" && c.Request.URL.Path == "/" {
			c.Next()
			return
		}

		// Apply auth to all other endpoints including metrics
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.Header("WWW-Authenticate", "Basic realm=\"coda\"")
			c.AbortWithStatus(401)
			return
		}
		if auth != "Basic "+*basicAuth {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	})
}
