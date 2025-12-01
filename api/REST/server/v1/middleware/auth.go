package middleware

import (
	"log"
	"net/http"
	"os"
	"ots/helper"
	"ots/mongo/dbops"
	"ots/tokenstructs"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthenticateAccessOf(of string, c *gin.Context) {
	header := c.GetHeader("Authorization")

	token := strings.TrimPrefix(header, "Bearer ")
	log.Printf("Token: %s", token)

	// decrypt the token
	payload := &tokenstructs.AccessToken{}
	var footer string
	if err := helper.Token.DecryptToken(token, payload, &footer); err != nil {
		log.Printf("error decrypting token: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, "token might be malformed")
		return
	}

	// check expiry
	exp := payload.Exp
	if time.Now().After(exp) {
		// token has expired
		c.AbortWithStatusJSON(http.StatusUnauthorized, "token has expired")
		return
	}

	// Get ID
	id := payload.Id

	switch strings.ToLower(of) {
	case "admin":
		if _, err := dbops.GetAdminBy("id", id); err != nil {
			// admin probably does not exist
			c.AbortWithStatusJSON(http.StatusUnauthorized, "error getting admin details")
			return
		}

		// Set context
		c.Set("admin", payload)

	case "resolver":
		if _, err := dbops.GetResolverBy("id", id); err != nil {
			// resolver probably does not exist
			c.AbortWithStatusJSON(http.StatusUnauthorized, "error getting resolver details")
			return
		}

		// Set context
		c.Set("resolver", payload)

	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, "unsupported access authentication")
		return
	}

	// Next
	c.Next()
}

func AuthenticateAdminAccess(c *gin.Context) {
	AuthenticateAccessOf("admin", c)
}

func AuthenticateResolverAccess(c *gin.Context) {
	AuthenticateAccessOf("resolver", c)
}

func AuthenticateAPIKey(c *gin.Context) {
	bearer := c.GetHeader("Authorization")

	apikey := strings.TrimPrefix(bearer, "Bearer ")

	myapikeychain := os.Getenv("API_KEY_CHAIN")
	myapikeys := strings.Split(myapikeychain, "-")

	if !slices.Contains(myapikeys, apikey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "a valid api key is needed to access the route")
		return
	}

	c.Next()
}
