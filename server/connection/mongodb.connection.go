package connection

import (
	"backend/server/env"
	"backend/server/structure"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct{}

var mongoDbConnected = false

func (ref MongoDB) Connect() (*mongo.Client, context.Context, error) {
	ctx := context.Background()
	uri := env.GetMongoUrl()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if !mongoDbConnected {
		log.Println("✅ MongoDB Connected")
		mongoDbConnected = true
	}
	return client, ctx, nil
}

func (ref MongoDB) CreateIndex(ctx context.Context, database *mongo.Database, collectionName string, listIndex []structure.IndexMongoDB) error {
	collection := database.Collection(collectionName)
	for _, index := range listIndex {
		opt := options.Index()
		opt = opt.SetUnique(index.Unique)
		opt = opt.SetName(index.Name)
		_, err := collection.Indexes().CreateOne(
			ctx,
			mongo.IndexModel{
				Keys:    index.Keys,
				Options: opt,
			},
		)
		if err != nil {
			log.Fatalf("Error on CreateIndex: %v", err)
			return err
		}
	}
	return nil
}
