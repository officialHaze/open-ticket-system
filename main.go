package main

import (
	"log"
	"os"
	"ots/helper"
	"ots/settings"

	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	settings.Generate()

	// Load the env file
	if err := godotenv.Load(settings.MySettings.Get_UseEnv()); err != nil {
		log.Fatalf("error loading ENV: %v", err)
	}
	log.Println("ENV file loaded successfully.")

	// Setup mgm
	err := mgm.SetDefaultConfig(
		nil,
		settings.MySettings.Get_DBName(),
		options.Client().ApplyURI(settings.MySettings.Get_MongoURL()))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("MGM setup complete.")

}

func main() {
	env := os.Getenv("ENV")
	log.Printf("OTS running in %s environment.", env)

	// DB index setup
	errs := helper.EnsureAllIndexes()
	log.Printf("Index creation errors: %v", errs)

	// Add initial admins
	ids := helper.AddInitialAdmins()
	log.Printf("Insterted Admin IDs: %v", ids)
}
