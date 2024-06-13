package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -> main collection
type User struct {
	ID        primitive.ObjectID  `json:"_id,omitempty"        bson:"_id,omitempty"`
	Name      string              `json:"name"                 bson:"name"`
	ImageURL  *string             `json:"image_url,omitempty"  bson:"image_url,omitempty"`
	Email     *string             `json:"email,omitempty"      bson:"email,omitempty"`
	Username  string              `json:"username"             bson:"username"`
	Password  string              `json:"password"             bson:"password"`
	IsVerify  bool                `json:"is_verify"            bson:"is_verify"` // jika login by email / baru, maka auto true
	IsActive  bool                `json:"is_active"            bson:"is_active"`
	RoleID    string              `json:"role_id"              bson:"role_id"`
	CreatedAt primitive.DateTime  `json:"created_at"           bson:"created_at"`
	UpdatedAt *primitive.DateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *primitive.DateTime `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// ------------------------------------------------------------------------

type UserRegisterBody struct {
	Name     string `json:"name" bson:"name"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type UserLoginBody struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type UserForgotPasswordBody struct {
	Email    *string `json:"email,omitempty"    bson:"email,omitempty"`
	Username *string `json:"username,omitempty" bson:"username,omitempty"`
}

type UserForgotPasswordResendBody struct {
	Ref string `json:"ref" bson:"ref"`
}

type UserForgotPasswordOtpValidBody struct {
	Ref  string `json:"ref"  bson:"ref"`
	Code string `json:"code" bson:"code"`
}

type UserForgotPasswordOtpSubmitBody struct {
	Ref      string `json:"ref"      bson:"ref"`
	Code     string `json:"code"     bson:"code"`
	Password string `json:"password" bson:"password"`
}
