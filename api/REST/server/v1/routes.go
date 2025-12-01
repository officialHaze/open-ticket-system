package v1

import (
	"ots/api/REST/server/v1/controller"
	"ots/api/REST/server/v1/middleware"
	"ots/ticketstructs"

	"github.com/gin-gonic/gin"
)

func MapRoutes(v1 *gin.RouterGroup) {
	v1.GET("/", controller.Home)

	// **** Admin Group **** //
	admin := v1.Group("/admin")
	admin.POST("/login", middleware.AuthenticateAPIKey, controller.AdminLogin)

	// **** Ticket Group **** //
	ticket := v1.Group("/ticket")
	ticket.POST("/new", middleware.AuthenticateAPIKey, controller.NewTicket)
	ticket.GET("/", middleware.AuthenticateAPIKey, controller.GetTicketsByCreator)

	ticket.PUT("/open",
		middleware.AuthenticateResolverAccess,
		func(ctx *gin.Context) { ctx.Set("ticketstatus", ticketstructs.GenerateTicketStatus().Open); ctx.Next() },
		controller.SetTicketStatus) // only resolvers can update status of tickets

	ticket.PUT("/inprogress",
		middleware.AuthenticateResolverAccess,
		func(ctx *gin.Context) {
			ctx.Set("ticketstatus", ticketstructs.GenerateTicketStatus().InProgress)
			ctx.Next()
		},
		controller.SetTicketStatus) // only resolvers can update status of tickets

	ticket.DELETE("/close",
		middleware.AuthenticateResolverAccess,
		func(ctx *gin.Context) {
			ctx.Set("ticketstatus", ticketstructs.GenerateTicketStatus().Closed)
			ctx.Next()
		},
		controller.SetTicketStatus) // only resolvers can update status of tickets

	// **** Resolver Group **** //
	resolver := v1.Group("/resolver")
	resolver.POST("/login", middleware.AuthenticateAPIKey, controller.ResolverLogin)
	resolver.POST("/new", middleware.AuthenticateAdminAccess, controller.NewResolver) // only admins can add new resolvers
	resolver.GET("/tickets", middleware.AuthenticateResolverAccess, controller.GetAssignedTickets)
}
