package routes

import (
	v1 "ots/api/REST/server/v1"
	"ots/api/REST/server/v1/middleware"

	"github.com/gin-gonic/gin"
)

func MapRoutes(r *gin.Engine) {
	api := r.Group("/api",middleware.RateLimmiter())

	// Version 1
	v1.MapRoutes(api.Group("/v1"))
}
