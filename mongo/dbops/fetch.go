package dbops

import (
	"fmt"
	"log"
	"ots/model"
	"ots/settings"
	"strings"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetSimilarTickets(title, description string) []*model.Ticket {
	similarTickets := make([]*model.Ticket, 0, 100)

	ticket := &model.Ticket{}
	coll := mgm.Coll(ticket)

	filter := bson.M{
		"$text": bson.M{
			"$search": title + " " + description,
		},
	}

	// Always add a similarity score to the document when searching and sort by the score
	// from most relevant to least relevant
	opts := options.Find()
	opts.SetProjection(bson.M{
		"similarityScore": bson.M{"$meta": "textScore"},
	})
	opts.SetSort(bson.M{
		"similarityScore": bson.M{"$meta": "textScore"},
	})

	cursor, err := coll.Find(settings.MySettings.Get_CtxWithTimeout(), filter, opts)
	if err != nil {
		log.Printf("error searching similar tickets based on TITLE - %s and DESCRIPTION - %s: %v", title, description, err)
		return []*model.Ticket{}
	}

	for cursor.Next(settings.MySettings.Get_CtxWithTimeout()) {
		var rec *model.Ticket
		if err := cursor.Decode(&rec); err != nil {
			log.Printf("error decoding document: %v", err)
			continue
		}

		similarTickets = append(similarTickets, rec)
	}

	// only first n (top) relevant documents
	n := min(5, len(similarTickets))

	return similarTickets[:n]
}

func GetAdminBy[T any](by string, d T) (*model.Admin, error) {
	admin := &model.Admin{}
	coll := mgm.Coll(admin)

	var filter bson.M

	switch strings.ToLower(by) {
	case "email":
		filter = bson.M{
			"email": d,
		}

	case "id":
		filter = bson.M{
			"_id": d,
		}

	default:
		return nil, fmt.Errorf("by factor not supported yet")
	}

	if err := coll.FindOne(settings.MySettings.Get_CtxWithTimeout(), filter).Decode(&admin); err != nil {
		// admin does not exist
		return nil, fmt.Errorf("admin does not exist")
	}

	return admin, nil
}

func GetResolverBy[T any](by string, d T) (*model.Resolver, error) {
	resolver := &model.Resolver{}
	coll := mgm.Coll(resolver)

	var filter bson.M

	switch strings.ToLower(by) {
	case "email":
		filter = bson.M{
			"email": d,
		}

	case "id":
		filter = bson.M{
			"_id": d,
		}

	default:
		return nil, fmt.Errorf("by factor not supported yet")
	}

	if err := coll.FindOne(settings.MySettings.Get_CtxWithTimeout(), filter).Decode(&resolver); err != nil {
		// admin does not exist
		return nil, fmt.Errorf("resolver does not exist")
	}

	return resolver, nil
}
