package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "You have reached the Home of V1 api")
}
