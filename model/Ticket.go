package model

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	mgm.DefaultModel `bson:",inline"`
	Title            string             `json:"title" bson:"title"`
	Description      string             `json:"description" bson:"description"`
	Status           string             `json:"status" bson:"status"`
	Priority         string             `json:"priority" bson:"priority"`
	Timeline         time.Time          `json:"timeline" bson:"timeline"`
	Milestones       []*TicketMilestone `json:"milestones" bson:"milestones"`
	AssignedTo       primitive.ObjectID `json:"assignedTo" bson:"assignedTo"`
	CreatorId        string             `json:"creatorId" bson:"creatorId"`
}

type TicketPipeline struct {
	mgm.DefaultModel `bson:",inline"`
	TicketID         primitive.ObjectID `json:"ticketId" bson:"ticketId"`
	ResolverID       primitive.ObjectID `json:"resolverId" bson:"resolverId"`
}

type TicketMilestone struct {
	Mark    int    `json:"mark" bson:"mark"`
	Title   string `json:"title" bson:"title"`
	Message string `json:"message" bson:"message"`
}
