package controller

import (
	"log"
	"net/http"
	"ots/model"
	"ots/mongo/dbops"
	"ots/pipeline"

	"github.com/gin-gonic/gin"
)

func NewTicket(c *gin.Context) {
	forcecreate := c.Query("force")

	ticketdetails := &model.Ticket{}
	if err := c.BindJSON(ticketdetails); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.Abort()
		return
	}

	if forcecreate == "0" {
		// First search and get similar tickets
		// if found, then return the similar tickets
		// no need to create new unless force query is true
		similarTickets := dbops.GetSimilarTickets(ticketdetails.Title, ticketdetails.Description)
		if len(similarTickets) > 0 {
			c.IndentedJSON(http.StatusAccepted, map[string]any{
				"message": "found similar tickets",
				"tickets": similarTickets,
			})
			return
		}
	}

	ticket, err := dbops.AddTicket(ticketdetails)
	if err != nil {
		log.Printf("error adding new ticket: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, "error adding new ticket. internal server error.")
		return
	}
	log.Printf("Ticket created with ID: %s", ticket.ID)

	c.IndentedJSON(http.StatusCreated, ticket.ID)

	// Push to ticket pipeline
	log.Printf("Pushing ticket with ID - %s to the pipeline.", ticket.ID)
	pipeline.TicketPipeline.Push(ticket)
	log.Printf("Ticket with ID - %s, pushed to pipeline.", ticket.ID)
}
