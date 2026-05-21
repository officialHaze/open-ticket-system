package helper

import (
	"fmt"
	"log"
	"os"
	"ots/mongo/dbops"
	"ots/settings"
	"strings"
)

func AddInitialAdmins() []interface{} {
	adminDetails := settings.MySettings.Get_InitialAdmins()

	insertedIds := make([]interface{}, 0, 100)

	for i, admin := range adminDetails {
		log.Printf("Admin email: %s", admin.Email)

		// Replace the password pointer with original password
		pass := os.Getenv(fmt.Sprintf("ADMIN_PASS_%d", i))
		admin.Password = strings.Replace(admin.Password, fmt.Sprintf("<ADMIN_PASSWORD_%d>", i), pass, 1)

		// Update/Set necessary fields
		admin.IsVerified = true
		admin.HasDefPassChanged = true

		// Hash and then save the password
		hashed, err := HashPasswd(admin.Password)
		if err != nil {
			log.Println(err)
			return []interface{}{}
		}
		admin.Password = hashed

		id, err := dbops.AddAdmin(admin)
		if err != nil {
			log.Printf("error adding admin with email - %s: %v", admin.Email, err)
			continue
		}

		insertedIds = append(insertedIds, id)
	}

	return insertedIds
}
