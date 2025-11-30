package server

import (
	"fmt"
	"log"
	"ots/api/REST/server/routes"
	"ots/settings"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	// Map routes
	routes.MapRoutes(r)

	addr := fmt.Sprintf(":%d", settings.MySettings.Get_ServerPort())
	if err := r.Run(addr); err != nil {
		log.Fatalln(err)
	}
}
