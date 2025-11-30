package controller

import (
	"log"
	"net/http"
	"ots/helper"
	"ots/model"
	"ots/mongo/dbops"

	"github.com/gin-gonic/gin"
)

func NewResolver(c *gin.Context) {
	resolvers := make([]*model.Resolver, 0, 100)
	if err := c.BindJSON(&resolvers); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.Abort()
		return
	}

	addedResolverIds := make([]string, len(resolvers))

	for i, r := range resolvers {
		// hash the password and update
		hashed, err := helper.HashPasswd(r.Password)
		if err != nil {
			log.Printf("error hashing password for resolver with email - %s", r.Email)
			continue
		}

		r.Password = hashed // replace plain text password with hash

		resolverId, err := dbops.AddResolver(r)
		if err != nil {
			log.Printf("error adding new resolver with email - %s: %v", r.Email, err)
			continue
		}
		log.Printf("Resolver added with ID: %s", resolverId.Hex())
		addedResolverIds[i] = resolverId.Hex()
	}

	c.IndentedJSON(http.StatusCreated, addedResolverIds)
}
