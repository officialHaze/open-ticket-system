package server

import (
	"fmt"
	"log"
	"ots/api/REST/server/routes"
	"ots/settings"
	"ots/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {
	if !util.InDevMode() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Cors setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     settings.MySettings.Get_AllowedClientOrigins(),
		AllowCredentials: true,
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Type"},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
	}))

	// Map routes
	routes.MapRoutes(r)

	addr := fmt.Sprintf(":%d", settings.MySettings.Get_ServerPort())
	if err := r.Run(addr); err != nil {
		log.Fatalln(err)
	}
}
