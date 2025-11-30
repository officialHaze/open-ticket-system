package dbops

import (
	"fmt"
	"log"
	"ots/model"
	"ots/settings"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAdmin(admin *model.Admin) (interface{}, error) {
	// defer settings.MySettings.Get_CtxCancel()()

	coll := mgm.Coll(admin)

	filter := bson.M{
		"email": admin.Email,
	}

	if err := coll.FindOne(settings.MySettings.Get_CtxWithTimeout(), filter).Decode(&admin); err != nil {
		// Admin does not exist
		// Insert one
		err := coll.Create(admin)
		if err != nil {
			return nil, err
		}
	}

	return admin.ID, nil
}

func AddResolver(resolver *model.Resolver) (primitive.ObjectID, error) {
	// defer settings.MySettings.Get_CtxCancel()()

	coll := mgm.Coll(resolver)

	filter := bson.M{
		"email": resolver.Email,
	}

	if err := coll.FindOne(settings.MySettings.Get_CtxWithTimeout(), filter).Decode(&resolver); err != nil {
		// Resolver does not exist
		// Insert one
		err := coll.Create(resolver)
		if err != nil {
			return primitive.NilObjectID, err
		}

		// Create a tracker for this resolver
		tt, err := AddTicketTracker(primitive.NilObjectID, resolver.ID)
		if err != nil {
			return primitive.NilObjectID, err
		}
		log.Printf("Ticket tracker created - %s, for resolver - %s", tt.ID, resolver.ID)

		// return resolver.ID, nil
	}

	return resolver.ID, nil
}

func AddTicket(ticket *model.Ticket) (*model.Ticket, error) {
	// defer settings.MySettings.Get_CtxCancel()()

	coll := mgm.Coll(ticket)

	filter := bson.M{
		"title": ticket.Title,
	}

	if err := coll.FindOne(settings.MySettings.Get_CtxWithTimeout(), filter).Decode(&ticket); err != nil {
		// ticket does not exist
		// Insert one
		ticket.AssignedTo = primitive.NilObjectID
		ticket.Milestones = append(ticket.Milestones, settings.MySettings.Get_DefaultTicketMilestones()[0])
		err := coll.Create(ticket)
		if err != nil {
			return nil, err
		}

		return ticket, nil
	}

	return nil, fmt.Errorf("ticket with duplicate title - %s, already exists", ticket.Title)
}

func AddTicketTracker(ticketId primitive.ObjectID, resolverId primitive.ObjectID) (*model.TicketTracker, error) {
	// defer settings.MySettings.Get_CtxCancel()()

	tickettracker := &model.TicketTracker{}
	coll := mgm.Coll(tickettracker)

	filter := bson.M{
		"$search": bson.M{
			"ticketId":   ticketId,
			"resolverId": resolverId,
		},
	}

	if err := coll.FindOne(settings.MySettings.Get_CtxWithTimeout(), filter).Decode(&tickettracker); err != nil {
		// ticket does not exist
		// Insert one
		tickettracker.TicketID = ticketId
		tickettracker.ResolverID = resolverId
		if err := coll.Create(tickettracker); err != nil {
			return nil, err
		}

		return tickettracker, nil
	}

	return nil, fmt.Errorf("ticket with ID - %s, is already assigned to resolver with ID - %s", ticketId, resolverId)
}
