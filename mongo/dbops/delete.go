package dbops

import (
	"context"
	"fmt"
	"ots/model"
	"ots/settings"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteTicketTracker(ticketId, resolverId primitive.ObjectID) error {
	ctxbase := context.TODO()
	ctx, cancel := context.WithTimeout(ctxbase, settings.MySettings.Get_CtxTimeout())
	defer cancel()

	tickettracker := &model.TicketTracker{}
	coll := mgm.Coll(tickettracker)

	filter := bson.M{
		"ticketId":   ticketId,
		"resolverId": resolverId,
	}

	if err := coll.FindOneAndDelete(ctx, filter).Decode(&tickettracker); err != nil {
		return fmt.Errorf("error deleting tracker record of Ticket ID - %s and Resolver ID - %s: %v", ticketId, resolverId, err)
	}

	return nil
}
