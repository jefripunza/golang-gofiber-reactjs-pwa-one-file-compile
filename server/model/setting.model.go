package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// -> main collection
type Setting struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Key   string             `json:"key"           bson:"key"`
	Value string             `json:"value"         bson:"value"`
}
