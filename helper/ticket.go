package helper

import (
	"log"
	"ots/pipeline"
	"ots/ticketassigner"
)

func GenerateTicketPipeline() {
	// Generate and build the middleware
	pipeline.GenerateTicketPipeline()
	log.Println("Ticket pipeline generated.")
	pipeline.TicketPipeline.Build(-1) // with default size
	log.Println("Ticket pipeline built.")
}

func InitializeTicketAssigner() {
	t := ticketassigner.New()
	t.Init()
}
