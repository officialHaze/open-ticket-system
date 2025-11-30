package controller

import (
	"log"
	"net/http"
	"ots/model"
	"ots/mongo/dbops"
	"ots/pipeline"
	"ots/ticketstructs"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetTicketsByCreator(c *gin.Context) {
	creatorid := c.Query("creatorid")

	tickets := dbops.GetTicketsBy("creatorid", creatorid)

	c.IndentedJSON(http.StatusOK, tickets)
}

func SetTicketStatusOpen(c *gin.Context) {
	ticketId := c.Query("ticketid")

	ticketObjectId, err := primitive.ObjectIDFromHex(ticketId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid ticket id")
		return
	}

	if err := dbops.UpdateTicketStatus(ticketstructs.GenerateTicketStatus().Open, ticketObjectId); err != nil {
		log.Printf("error opening ticket: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "error opening ticket. internal server error")
		return
	}

	c.IndentedJSON(http.StatusCreated, "ticket has been opened")
}
