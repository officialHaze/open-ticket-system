package dbops

import (
	"fmt"
	"ots/model"
	"ots/settings"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAdmin(admin *model.Admin) (interface{}, error) {
	defer settings.MySettings.Get_CtxCancel()()

	coll := mgm.Coll(admin)

	filter := bson.M{
		"$text": bson.M{
			"$search": admin.Email,
		},
	}

	if err := coll.FindOne(settings.MySettings.Get_CtxWithTimeout(), filter).Decode(&admin); err != nil {
		// Admin does not exist
		// Insert one
		err := coll.Create(admin)
		if err != nil {
			return nil, err
		}

		return admin.ID, nil
	}

	return nil, fmt.Errorf("admin with email - %s, already exists", admin.Email)
}

func AddResolver(resolver *model.Resolver) (primitive.ObjectID, error) {
	defer settings.MySettings.Get_CtxCancel()()

	coll := mgm.Coll(resolver)

	filter := bson.M{
		"$text": bson.M{
			"$search": resolver.Email,
		},
	}

	if err := coll.FindOne(settings.MySettings.Get_CtxWithTimeout(), filter).Decode(&resolver); err != nil {
		// Resolver does not exist
		// Insert one
		err := coll.Create(resolver)
		if err != nil {
			return primitive.NilObjectID, err
		}

		return resolver.ID, nil
	}

	return primitive.NilObjectID, fmt.Errorf("resolver with email - %s, already exists", resolver.Email)
}
