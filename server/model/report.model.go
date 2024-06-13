package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -> main collection
type Report struct {
	ID        primitive.ObjectID  `json:"_id,omitempty"        bson:"_id,omitempty"`
	Username  string              `json:"username"             bson:"username" unique:"true"`
	Password  string              `json:"password"             bson:"password"`
	Name      string              `json:"name"                 bson:"name"`
	ImageURL  string              `json:"image_url"            bson:"image_url"`
	CreatedAt primitive.DateTime  `json:"created_at"           bson:"created_at"`
	UpdatedAt *primitive.DateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *primitive.DateTime `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
