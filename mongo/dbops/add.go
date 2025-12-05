package dbops

import (
	"context"
	"fmt"
	"log"
	"ots/model"
	"ots/settings"
	"ots/ticketstructs"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAdmin(admin *model.Admin) (interface{}, error) {
	ctxbase := context.TODO()
	ctx, cancel := context.WithTimeout(ctxbase, settings.MySettings.Get_CtxTimeout())
	defer cancel()

	coll := mgm.Coll(admin)

	filter := bson.M{
		"email": admin.Email,
	}

	if err := coll.FindOne(ctx, filter).Decode(&admin); err != nil {
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
	ctxbase := context.TODO()
	ctx, cancel := context.WithTimeout(ctxbase, settings.MySettings.Get_CtxTimeout())
	defer cancel()

	coll := mgm.Coll(resolver)

	filter := bson.M{
		"email": resolver.Email,
	}

	if err := coll.FindOne(ctx, filter).Decode(&resolver); err != nil {
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
	ctxbase := context.TODO()
	ctx, cancel := context.WithTimeout(ctxbase, settings.MySettings.Get_CtxTimeout())
	defer cancel()

	coll := mgm.Coll(ticket)

	filter := bson.M{
		"title": ticket.Title,
	}

	if err := coll.FindOne(ctx, filter).Decode(&ticket); err != nil {
		// ticket does not exist
		// Insert one
		ticket.Milestones = []*model.TicketMilestone{} // set empty milestone slice during creation
		err := coll.Create(ticket)
		if err != nil {
			return nil, err
		}

		// Update ticket status to created
		err = UpdateTicketStatus(ticketstructs.GenerateTicketStatus().Created, ticket.ID)
		if err != nil {
			log.Printf("error updating ticket status to created for - %s: %v", ticket.ID, err)
		}

		// Assign nil resolver to ticket in the begining
		_, err = AssignResolverToTicket(ticket, primitive.NilObjectID)
		if err != nil {
			log.Printf("error assigning nil resolver ID to ticket - %s: %v", ticket.ID, err)
		}

		// Append the default begining milestone
		err = AppendTicketMileStone(settings.MySettings.Get_DefaultTicketMilestones()[0], ticket.ID)
		if err != nil {
			log.Printf("error appending milestone to ticket - %s: %v", ticket.ID, err)
		}

		return ticket, nil
	}

	return nil, fmt.Errorf("ticket with duplicate title - %s, already exists", ticket.Title)
}

func AddTicketTracker(ticketId primitive.ObjectID, resolverId primitive.ObjectID) (*model.TicketTracker, error) {
	ctxbase := context.TODO()
	ctx, cancel := context.WithTimeout(ctxbase, settings.MySettings.Get_CtxTimeout())
	defer cancel()

	tickettracker := &model.TicketTracker{}
	coll := mgm.Coll(tickettracker)

	filter := bson.M{
		"$search": bson.M{
			"ticketId":   ticketId,
			"resolverId": resolverId,
		},
	}

	if err := coll.FindOne(ctx, filter).Decode(&tickettracker); err != nil {
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
