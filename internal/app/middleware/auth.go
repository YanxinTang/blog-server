package middleware

import (
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("login") != true {
			c.Error(e.ERROR_SESSION_EXPIRED)
			c.Abort()
		}
		c.Next()
	}
}
