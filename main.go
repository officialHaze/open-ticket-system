package main

import (
	"log"
	"os"
	"ots/api/REST/server"
	"ots/helper"
	"ots/settings"
	"strings"

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

	switch getArg() {
	case "ots-server":
		helper.GenerateTicketPipeline()
		go func() {
			helper.InitializeTicketAssigner()
		}()
		server.Start()
		return
	default:
		log.Fatalln("Unsupported arg")
		return
	}
}

func getArg() string {
	args := os.Args

	execname := args[0]
	log.Printf("Executable name: %s", execname)

	if len(args) <= 1 {
		log.Fatalln("No arg provided.")
	}

	mainarg := args[1]
	mainarg = strings.TrimSpace(strings.ToLower(mainarg)) // normalize
	log.Printf("Main Arg: %s", mainarg)

	return mainarg
}
