package Models

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Zila struct {
	ID            primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string               `json:"name,omitempty"`
	Established	  string               `json:"established,omitempty"`
	Status        bool                 `json:"status,omitempty"`
	Longitude     float64              `json:"longitude,omitempty"`
	Latitude      float64              `json:"latitude,omitempty"`
	ParentId 	  primitive.ObjectID   `json:"parentid,omitempty"`
	Upazilas 	  []primitive.ObjectID `json:"upazilas,omitempty"`
}

func (obj Zila) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.Established, validation.Required), 
		validation.Field(&obj.Longitude, validation.Required),
		validation.Field(&obj.Latitude, validation.Required),
		validation.Field(&obj.ParentId, validation.Required),
	)
}

type ZilaSearch struct {
	ID           		primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	IDIsUsed     		bool               `json:"idisused,omitempty"`
	ParentId     		primitive.ObjectID `json:"parentid"`
	ParentIdIsUsed      bool               `json:"parentidisused,omitempty"`
	Name         		string             `json:"name,omitempty"`
	NameIsUsed   		bool               `json:"nameisused,omitempty"`
	Status       		bool               `json:"status,omitempty"`
	StatusIsUsed 		bool               `json:"statusisused,omitempty"`
}

func (obj ZilaSearch) GetZilaSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}
	if obj.ParentIdIsUsed {
		self["parentid"] = obj.ParentId
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	return self
}