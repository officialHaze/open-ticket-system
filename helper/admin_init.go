package helper

import (
	"log"
	"ots/mongo/dbops"
	"ots/settings"
)

func AddInitialAdmins() []interface{} {
	adminDetails := settings.MySettings.Get_InitialAdmins()

	insertedIds := make([]interface{}, 0, 100)

	for _, admin := range adminDetails {
		log.Printf("Admin email: %s", admin.Email)

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
