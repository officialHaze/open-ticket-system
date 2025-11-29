package dbops

import (
	"fmt"
	"ots/model"
	"ots/settings"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AssignResolverToTicket(ticket *model.Ticket, resolverId primitive.ObjectID) (*model.Ticket, error) {
	// defer settings.MySettings.Get_CtxCancel()()

	coll := mgm.Coll(ticket)

	filter := bson.M{
		"_id": ticket.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"assignedTo": resolverId,
		},
	}
	if err := coll.FindOneAndUpdate(settings.MySettings.Get_CtxWithTimeout(), filter, update).Decode(&ticket); err != nil {
		// ticket does not exist
		return nil, fmt.Errorf("ticket with ID - %s does not exist: %v", ticket.ID, err)
	}

	return ticket, nil
}

// func UpdateTicketTracker(resolverId primitive.ObjectID)
