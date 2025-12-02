package dbops

import (
	"context"
	"fmt"
	"ots/model"
	"ots/settings"
	"ots/ticketstructs"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AssignResolverToTicket(ticket *model.Ticket, resolverId primitive.ObjectID) (*model.Ticket, error) {
	ctxbase := context.TODO()
	ctx, cancel := context.WithTimeout(ctxbase, settings.MySettings.Get_CtxTimeout())
	defer cancel()

	coll := mgm.Coll(ticket)

	filter := bson.M{
		"_id": ticket.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"assignedTo": resolverId,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&ticket); err != nil {
		// ticket does not exist
		return nil, fmt.Errorf("ticket with ID - %s does not exist: %v", ticket.ID, err)
	}

	return ticket, nil
}

func UpdateTicketStatus(status string, ticketId primitive.ObjectID) error {
	ctxbase := context.TODO()
	ctx, cancel := context.WithTimeout(ctxbase, settings.MySettings.Get_CtxTimeout())
	defer cancel()

	ts := ticketstructs.GenerateTicketStatus()
	if !ts.IsValidStatus(status) {
		return fmt.Errorf("invalid ticket status to update")
	}

	ticket := &model.Ticket{}
	coll := mgm.Coll(ticket)

	searchfilter := bson.M{
		"_id": ticketId,
	}

	updatefilter := bson.M{
		"$set": bson.M{
			"status":           status,
			"statusUpdateDate": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := coll.FindOneAndUpdate(ctx, searchfilter, updatefilter, opts).Decode(&ticket)
	if err != nil {
		return fmt.Errorf("error updating ticket status")
	}

	return nil
}

func AppendTicketMileStone(milestone *model.TicketMilestone, ticketId primitive.ObjectID) error {
	ctxbase := context.TODO()
	ctx, cancel := context.WithTimeout(ctxbase, settings.MySettings.Get_CtxTimeout())
	defer cancel()

	ticket := &model.Ticket{}
	coll := mgm.Coll(ticket)

	searchExpr := bson.M{
		"_id": ticketId,
	}

	if err := coll.FindOne(ctx, searchExpr).Decode(&ticket); err != nil {
		return fmt.Errorf("ticket with ID - %s does not exist: %v", ticketId, err)
	}

	existingMilestones := ticket.Milestones

	for _, exmilestone := range existingMilestones {
		if exmilestone.Mark == milestone.Mark {
			// milestone already exists
			return nil
		}
	}

	updateExpr := bson.M{
		"$push": bson.M{
			"milestones": milestone,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := coll.FindOneAndUpdate(ctx, searchExpr, updateExpr, opts).Decode(&ticket); err != nil {
		return fmt.Errorf("error appending milestone to ticket: %v", err)
	}

	return nil
}

func SetPriority(priority string, ticketId primitive.ObjectID) error {
	ctxbase := context.TODO()
	ctx, cancel := context.WithTimeout(ctxbase, settings.MySettings.Get_CtxTimeout())
	defer cancel()

	ticket := &model.Ticket{}
	coll := mgm.Coll(ticket)

	searchfilter := bson.M{
		"_id": ticketId,
	}

	updatefilter := bson.M{
		"$set": bson.M{
			"priority": priority,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := coll.FindOneAndUpdate(ctx, searchfilter, updatefilter, opts).Decode(&ticket)
	if err != nil {
		return fmt.Errorf("error setting priority for ticket - %s", ticketId)
	}

	return nil
}
