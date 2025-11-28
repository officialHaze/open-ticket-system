package main

import (
	"log"
	"os"
	"ots/util"
)

func init() {
	// Load the env file
	if err := util.LoadEnv(); err != nil {
		log.Fatalln(err)
	}
	log.Println("ENV file loaded successfully.")
}

func main() {
	env := os.Getenv("ENV")
	log.Printf("OTS running in %s environment.", env)
}
