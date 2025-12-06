package controller

import (
	"fmt"
	"log"
	"net/http"
	"ots/model"
	"ots/mongo/dbops"
	"ots/pipeline"
	"ots/settings"
	"ots/ticketstructs"
	"ots/tokenstructs"

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

	if ticketdetails.CreatorId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Mandatory Creator ID is missing")
		return
	}

	if forcecreate == "0" {
		// First search and get similar tickets
		// if found, then return the similar tickets
		// no need to create new unless force query is true
		similarTickets := dbops.GetSimilarTickets(ticketdetails.Title, ticketdetails.Description, ticketdetails.CreatorId)
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

func SetTicketStatus(c *gin.Context) {
	resolverpayload, exists := c.Get("resolver")
	if !exists {
		c.IndentedJSON(http.StatusForbidden, "resolver session not found")
		return
	}

	resolverId := resolverpayload.(*tokenstructs.AccessToken).Id

	ticketstatus := c.GetString("ticketstatus")
	ticketId := c.Query("ticketid")

	ticketObjectId, err := primitive.ObjectIDFromHex(ticketId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid ticket id")
		return
	}

	// Pre-update Additional methods
	switch ticketstatus {
	case ticketstructs.GenerateTicketStatus().Open:
		// When a ticket is opened, append appropriate milestone
		dbops.AppendTicketMileStone(settings.MySettings.Get_DefaultTicketMilestones()[1], ticketObjectId)

	case ticketstructs.GenerateTicketStatus().InProgress:
		// When a ticket is in progress, append appropriate milestone
		dbops.AppendTicketMileStone(settings.MySettings.Get_DefaultTicketMilestones()[2], ticketObjectId)

	case ticketstructs.GenerateTicketStatus().Closed:
		// When a ticket is closed
		// It should be removed from the tracker
		// of the resolver
		if err := dbops.DeleteTicketTracker(ticketObjectId, resolverId); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "error updating ticket status. the ticket might already be closed. internal error")
			return
		}
		// appr milestone should be added
		dbops.AppendTicketMileStone(settings.MySettings.Get_DefaultTicketMilestones()[3], ticketObjectId)

	default:
		return
	}

	if err := dbops.UpdateTicketStatus(ticketstatus, ticketObjectId); err != nil {
		log.Printf("error updating ticket status: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "error updating ticket status. internal server error")
		return
	}

	c.IndentedJSON(http.StatusCreated, "ticket status updated")
}

func SetTicketPriority(c *gin.Context) {
	ticketId := c.Query("ticketid")

	priority := c.Param("set")
	p := &ticketstructs.Priority{}

	if !p.IsValid(priority) {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid priority type provided")
		return
	}

	ticketObjectId, err := primitive.ObjectIDFromHex(ticketId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid ticket id")
		return
	}

	if err := dbops.SetPriority(priority, ticketObjectId); err != nil {
		log.Printf("error setting priority: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "error setting ticket priority. internal server error")
		return
	}

	c.IndentedJSON(http.StatusCreated, fmt.Sprintf("priority set to - %s for ticket - %s", priority, ticketId))
}
