package model

import "github.com/kamva/mgm/v3"

type Admin struct {
	mgm.DefaultModel  `bson:",inline"`
	Name              string `json:"name" bson:"name"`
	Email             string `json:"email" bson:"email"` // Unique
	Phone             string `json:"phone" bson:"phone"`
	Password          string `json:"password" bson:"password"`
	IsVerified        bool   `json:"isVerified" bson:"isVerified"`
	HasDefPassChanged bool   `json:"hasDefPassChanged" bson:"hasDefPassChanged"`
}
