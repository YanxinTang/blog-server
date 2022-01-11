package middleware

import (
	"log"
	"net/http"

	"github.com/YanxinTang/blog-server/e"
	"github.com/gin-gonic/gin"
)

type BadRequestMeta struct {
	Url string
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			lastError := c.Errors.Last().Err
			switch err := lastError.(type) {
			case e.ApiError:
				c.AbortWithStatusJSON(err.Code, gin.H{
					"message": err.Message,
				})

			default:
				log.Println(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "内部错误",
				})
			}
		}
	}
}
