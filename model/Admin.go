package model

import "github.com/kamva/mgm/v3"

type Admin struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Email            string `json:"email" bson:"email"`
	Phone            string `json:"phone" bson:"phone"`
}
