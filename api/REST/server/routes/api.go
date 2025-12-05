package routes

import (
	v1 "ots/api/REST/server/v1"

	"github.com/gin-gonic/gin"
)

func MapRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Version 1
	v1.MapRoutes(api.Group("/v1"))
}
