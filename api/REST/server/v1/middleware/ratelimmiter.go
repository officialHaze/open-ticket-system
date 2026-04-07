package middleware

import (
	"ots/util"
	"net/http"
	"github.com/gin-gonic/gin"
)

func RateLimmiter() gin.HandlerFunc{
	return func(c *gin.Context) {
		select{
		case <-util.Limmiter:
			c.Next()
		default:
			//if no token return the error message and abort the request
			c.JSON(http.StatusTooManyRequests,gin.H{"error":"Too many Request"})
			c.Abort()
		}
	}
}