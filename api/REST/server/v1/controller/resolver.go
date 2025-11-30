package controller

import (
	"log"
	"net/http"
	"ots/model"
	"ots/mongo/dbops"

	"github.com/gin-gonic/gin"
)

func NewResolver(c *gin.Context) {
	resolverDetails := &model.Resolver{}
	if err := c.BindJSON(resolverDetails); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.Abort()
		return
	}

	resolverId, err := dbops.AddResolver(resolverDetails)
	if err != nil {
		log.Printf("error adding new resolver: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, "error adding new resolver. internal server error.")
		return
	}
	log.Printf("Resolver added with ID: %s", resolverId.Hex())

	c.IndentedJSON(http.StatusCreated, resolverId.Hex())
}
