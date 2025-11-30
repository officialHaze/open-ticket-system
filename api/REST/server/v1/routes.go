package v1

import (
	"ots/api/REST/server/v1/controller"
	"ots/api/REST/server/v1/middleware"

	"github.com/gin-gonic/gin"
)

func MapRoutes(v1 *gin.RouterGroup) {
	v1.GET("/", controller.Home)

	// **** Admin Group **** //
	admin := v1.Group("/admin")
	admin.POST("/login", controller.AdminLogin)

	// **** Ticket Group **** //
	ticket := v1.Group("/ticket")
	ticket.POST("/new", controller.NewTicket)

	// **** Resolver Group **** //
	resolver := v1.Group("/resolver")
	resolver.POST("/login", controller.ResolverLogin)
	resolver.POST("/new", middleware.AuthenticateAdminAccess, controller.NewResolver) // only admins can add new resolvers
}
