package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// -> main collection
type RevokeToken struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"user_id"       bson:"user_id"`
	JwtID     string             `json:"jwt_id"        bson:"jwt_id"`
	BrowserID string             `json:"browser_id"    bson:"browser_id"`
	ExpiredAt primitive.DateTime `json:"expired_at"    bson:"expired_at"`
	LoginAt   primitive.DateTime `json:"login_at"      bson:"login_at"`
}
