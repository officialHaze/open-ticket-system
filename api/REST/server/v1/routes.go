package v1

import (
	"ots/api/REST/server/v1/controller"

	"github.com/gin-gonic/gin"
)

func MapRoutes(v1 *gin.RouterGroup) {
	v1.GET("/", controller.Home)

	ticket := v1.Group("/ticket")
	ticket.POST("/new", controller.NewTicket)

	resolver := v1.Group("/resolver")
	resolver.POST("/new", controller.NewResolver)
}
