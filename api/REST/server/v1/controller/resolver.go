package controller

import (
	"fmt"
	"log"
	"net/http"
	"ots/helper"
	"ots/model"
	"ots/mongo/dbops"
	"ots/settings"
	"ots/tokenstructs"
	"time"

	"github.com/gin-gonic/gin"
)

func NewResolver(c *gin.Context) {
	resolvers := make([]*model.Resolver, 0, 100)
	if err := c.BindJSON(&resolvers); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.Abort()
		return
	}

	addedResolverIds := make([]string, 0, len(resolvers))
	invalidEmailIds, invalidPhonenos := make([]string, 0, len(resolvers)), make([]string, 0, len(resolvers))

	for _, r := range resolvers {
		// validate email
		if err := helper.ValidateEmail(r.Email); err != nil {
			invalidEmailIds = append(invalidEmailIds, fmt.Errorf("validation error for email - %s: %v", r.Email, err).Error())
			continue
		}

		// validate phone
		if err := helper.ValidatePhone(r.Phone); err != nil {
			invalidPhonenos = append(invalidPhonenos, fmt.Errorf("validation error for email - %s: %v", r.Email, err).Error())
			continue
		}

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
		addedResolverIds = append(addedResolverIds, resolverId.Hex())
	}

	c.IndentedJSON(http.StatusAccepted, map[string]any{
		"addedResolvers":  addedResolverIds,
		"invalidEmails":   invalidEmailIds,
		"invalidPhonenos": invalidPhonenos,
	})
}

func ResolverLogin(c *gin.Context) {
	resolvercreds := &model.Resolver{}
	if err := c.BindJSON(resolvercreds); err != nil {
		log.Printf("error while binding resolver creds to JSON: %v", err)
		c.Abort()
		return
	}

	email := resolvercreds.Email
	password := resolvercreds.Password

	// Get resolver details by email
	resolver, err := dbops.GetResolverBy("email", email)
	if err != nil {
		log.Printf("error getting resolver by email - %s: %v", email, err)
		c.IndentedJSON(http.StatusInternalServerError, "error getting resolver details")
		return
	}

	// Match the password
	if err := helper.CompareHash(resolver.Password, password); err != nil {
		log.Println(err)
		// passwords dont match
		c.IndentedJSON(http.StatusUnauthorized, "incorrect password")
		return
	}

	// Create access token
	tokenpayload := &tokenstructs.AccessToken{
		Id:    resolver.ID,
		Name:  resolver.Name,
		Email: resolver.Email,
		Exp:   time.Now().Add(settings.MySettings.Get_AccessTokenExpMin()),
	}
	accesstoken, err := helper.Token.CreateToken(tokenpayload)
	if err != nil {
		log.Printf("error generating access token for resolver (%s): %v", resolver.Email, err)
		c.IndentedJSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.IndentedJSON(http.StatusAccepted, map[string]string{
		"token": accesstoken,
		"name":  resolver.Name,
		"email": resolver.Email,
		"phone": resolver.Phone,
	})
}

func GetAssignedTickets(c *gin.Context) {
	resolverctx, exists := c.Get("resolver")
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, "resolver session unavailable")
		return
	}
	resolver := resolverctx.(*tokenstructs.AccessToken)
	log.Printf("Resolver ID: %s", resolver.Id)

	tickets := dbops.GetTicketsBy("assignee", resolver.Id)

	c.IndentedJSON(http.StatusOK, tickets)
}
