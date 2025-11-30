package middleware

import (
	"log"
	"net/http"
	"ots/helper"
	"ots/mongo/dbops"
	"ots/tokenstructs"
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
