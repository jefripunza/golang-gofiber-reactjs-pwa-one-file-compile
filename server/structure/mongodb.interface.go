package structure

import "go.mongodb.org/mongo-driver/bson"

type IndexMongoDB struct {
	Name   string `json:"name"`
	Unique bool   `json:"unique"`
	Keys   bson.D `json:"keys"`
}
