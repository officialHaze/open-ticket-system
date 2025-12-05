package controller

import (
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

func AdminLogin(c *gin.Context) {
	admincreds := &model.Admin{}
	if err := c.BindJSON(admincreds); err != nil {
		log.Printf("error while binding admin creds to JSON: %v", err)
		c.Abort()
		return
	}

	email := admincreds.Email
	password := admincreds.Password

	// Get admin details by email
	admin, err := dbops.GetAdminBy("email", email)
	if err != nil {
		log.Printf("error getting admin by email - %s: %v", email, err)
		c.IndentedJSON(http.StatusInternalServerError, "error getting admin details")
		return
	}

	// Match the password
	if err := helper.CompareHash(admin.Password, password); err != nil {
		log.Println(err)
		// passwords dont match
		c.IndentedJSON(http.StatusUnauthorized, "incorrect password")
		return
	}

	// Create access token
	tokenpayload := &tokenstructs.AccessToken{
		Id:    admin.ID,
		Name:  admin.Name,
		Email: admin.Email,
		Exp:   time.Now().Add(settings.MySettings.Get_AccessTokenExpMin()),
	}
	accesstoken, err := helper.Token.CreateToken(tokenpayload)
	if err != nil {
		log.Printf("error generating access token for admin (%s): %v", admin.Email, err)
		c.IndentedJSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.IndentedJSON(http.StatusAccepted, map[string]string{
		"token": accesstoken,
		"name":  admin.Name,
		"email": admin.Email,
		"phone": admin.Phone,
	})
}
