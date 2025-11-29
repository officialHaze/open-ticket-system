package dbops

import (
	"ots/model"
	"ots/settings"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTicketTrackers() []*model.TicketTracker {
	// defer settings.MySettings.Get_CtxCancel()

	trackers := make([]*model.TicketTracker, 0, 100)

	ticketTracker := &model.TicketTracker{}
	coll := mgm.Coll(ticketTracker)

	cursor, err := coll.Find(settings.MySettings.Get_CtxWithTimeout(), bson.M{})
	if err != nil {
		return []*model.TicketTracker{}
	}

	for cursor.Next(settings.MySettings.Get_CtxWithTimeout()) {
		var tracker *model.TicketTracker
		if err := cursor.Decode(&tracker); err != nil {
			continue
		}

		trackers = append(trackers, tracker)
	}

	return trackers
}
