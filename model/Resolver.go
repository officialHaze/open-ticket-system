package model

import "github.com/kamva/mgm/v3"

type Resolver struct {
	mgm.DefaultModel  `bson:",inline"`
	Name              string `json:"name" bson:"name"`
	Email             string `json:"email" bson:"email"` // unique
	Phone             string `json:"phone" bson:"phone"`
	Status            string `json:"status" bson:"status"`
	Password          string `json:"password" bson:"password"`
	IsVerified        bool   `json:"isVerified" bson:"isVerified"`
	HasDefPassChanged bool   `json:"hasDefPassChanged" bson:"hasDefPassChanged"`
}
