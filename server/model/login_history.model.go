package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// -> main collection
type LoginHistory struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"user_id"       bson:"user_id"`
	UserAgent string             `json:"user_agent"    bson:"user_agent"`
	LoginAt   primitive.DateTime `json:"login_at"      bson:"login_at"`
}
