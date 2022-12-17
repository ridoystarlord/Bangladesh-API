package Models

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Division struct {
	ID            primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string               `json:"name,omitempty"`
	IsoCode		  string               `json:"isocode,omitempty"`
	Established	  string               `json:"established,omitempty"`
	Status        bool                 `json:"status,omitempty"`
	Longitude     float64              `json:"longitude,omitempty"`
	Latitude      float64              `json:"latitude,omitempty"`
	Zilas 	  	  []primitive.ObjectID `json:"zilas,omitempty"`
}

func (obj Division) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.IsoCode, validation.Required),
		validation.Field(&obj.Established, validation.Required), 
		validation.Field(&obj.Longitude, validation.Required),
		validation.Field(&obj.Latitude, validation.Required),
	)
}

type DivisionSearch struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	IDIsUsed     bool               `json:"idisused,omitempty"`
	Name         string             `json:"name,omitempty"`
	NameIsUsed   bool               `json:"nameisused,omitempty"`
	Status       bool               `json:"status,omitempty"`
	StatusIsUsed bool               `json:"statusisused,omitempty"`
}

func (obj DivisionSearch) GetDivisionSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
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

type DivisionPopulated struct {
	ID            primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string               `json:"name,omitempty"`
	IsoCode		  string               `json:"isocode,omitempty"`
	Established	  string               `json:"established,omitempty"`
	Status        bool                 `json:"status,omitempty"`
	Longitude     float64              `json:"longitude,omitempty"`
	Latitude      float64              `json:"latitude,omitempty"`
	Zilas 	  	  []Zila `json:"zilas,omitempty"`
}

func (obj *DivisionPopulated) CloneFrom(other Division) {
	obj.ID = other.ID
	obj.Name = other.Name
	obj.IsoCode = other.IsoCode
	obj.Established = other.Established
	obj.Status = other.Status
	obj.Longitude = other.Longitude
	obj.Latitude = other.Latitude
	obj.Zilas = []Zila{}
}