package helper

import (
	"fmt"
	"log"
	"ots/settings"

	"golang.org/x/crypto/bcrypt"
)

func HashPasswd(plain string) (string, error) {
	if plain == "" {
		return "", fmt.Errorf("plan text is empty")
	}

	hashrounds := settings.MySettings.Get_PasswdHashRounds()

	plainb := []byte(plain)

	passb, err := bcrypt.GenerateFromPassword(plainb, hashrounds)
	if err != nil {
		return "", fmt.Errorf("error hashing password - %v", err)
	}
	log.Println("password hashed.")

	return string(passb), nil
}
