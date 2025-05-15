package middleware

import (
	"github.com/gin-gonic/gin"
)

func Version(r *gin.Engine, version string) {
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("x-coda-version", "coda/v"+version)
		c.Next()
	})
}
