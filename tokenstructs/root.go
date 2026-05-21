package tokenstructs

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessToken struct {
	Id    primitive.ObjectID
	Name  string
	Email string
	Exp   time.Time
}
