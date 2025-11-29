package mongo

import (
	"fmt"
	"log"
	"ots/model"
	"ots/settings"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ticket Index
func EnsureTicketIndexes() error {
	ticket := &model.Ticket{}
	coll := mgm.Coll(ticket)

	// index 1 (compound text index)
	index1 := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "title",
				Value: "text",
			},
			{
				Key:   "description",
				Value: "text",
			},
			{
				Key:   "status",
				Value: "text",
			},
			{
				Key:   "priority",
				Value: "text",
			},
		},
		Options: options.Index().SetName("text_compund_idx"),
	}

	// index 2
	index2 := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "timeline",
				Value: 1,
			},
		},
		Options: options.Index().SetName("timeline_idx"),
	}

	// index 3
	index3 := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "assignedTo",
				Value: 1,
			},
		},
		Options: options.Index().SetName("resolver_id_idx"),
	}

	// index 4
	index4 := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "creatorId",
				Value: 1,
			},
		},
		Options: options.Index().SetName("creator_id_idx"),
	}

	index5 := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "title",
				Value: 1,
			},
		},
		Options: options.Index().SetName("ticket_title_idx").SetUnique(true),
	}

	// create indexes
	idxs := []mongo.IndexModel{
		index1,
		index2,
		index3,
		index4,
		index5,
	}
	names, err := coll.Indexes().CreateMany(settings.MySettings.Get_CtxWithTimeout(), idxs)
	if err != nil {
		return fmt.Errorf("error creating all the indexes %v: %v", idxs, err)
	}

	log.Printf("Indexes created with names: %v", names)

	return nil
}

// Ticket Pipeline Index
func EnsureTicketTrackerIndexes() error {
	ticketpipeline := &model.TicketTracker{}
	coll := mgm.Coll(ticketpipeline)

	// index 1
	index1 := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "ticketId",
				Value: 1,
			},
		},
		Options: options.Index().SetName("ticket_id_idx"),
	}

	// index 2
	index2 := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "resolverId",
				Value: 1,
			},
		},
		Options: options.Index().SetName("resolver_id_idx"),
	}

	// create indexes
	idxs := []mongo.IndexModel{
		index1,
		index2,
	}
	names, err := coll.Indexes().CreateMany(settings.MySettings.Get_CtxWithTimeout(), idxs)
	if err != nil {
		return fmt.Errorf("error creating all the indexes %v: %v", idxs, err)
	}

	log.Printf("Indexes created with names: %v", names)

	return nil
}

// Resolver Index
func EnsureResolverIndexes() error {
	resolver := &model.Resolver{}
	coll := mgm.Coll(resolver)

	// index 1 (compound text index)
	index1 := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "name",
				Value: "text",
			},
			{
				Key:   "status",
				Value: "text",
			},
		},
		Options: options.Index().SetName("text_compund_idx"),
	}

	index2 := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "email",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true).SetName("resolver_email_idx"),
	}

	// create indexes
	idxs := []mongo.IndexModel{
		index1,
		index2,
	}
	names, err := coll.Indexes().CreateMany(settings.MySettings.Get_CtxWithTimeout(), idxs)
	if err != nil {
		return fmt.Errorf("error creating all the indexes %v: %v", idxs, err)
	}

	log.Printf("Indexes created with names: %v", names)

	return nil
}

// Admin Index
func EnsureAdminIndexes() error {
	admin := &model.Admin{}
	coll := mgm.Coll(admin)

	// index 1 (compound text index)
	index1 := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "name",
				Value: "text",
			},
		},
		Options: options.Index().SetName("text_compund_idx"),
	}

	// index 2
	index2 := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "email",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true).SetName("admin_email_index"),
	}

	// create indexes
	idxs := []mongo.IndexModel{
		index1,
		index2,
	}
	names, err := coll.Indexes().CreateMany(settings.MySettings.Get_CtxWithTimeout(), idxs)
	if err != nil {
		return fmt.Errorf("error creating all the indexes %v: %v", idxs, err)
	}

	log.Printf("Indexes created with names: %v", names)

	return nil
}
