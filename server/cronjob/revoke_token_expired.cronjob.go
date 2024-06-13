package cronjob

import (
	"backend/server/connection"
	"backend/server/env"
	"backend/server/variable"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func RevokeTokenExpired() {
	var err error
	ctx := context.Background()

	currentTime := time.Now()

	MongoDB := connection.MongoDB{}

	client, ctx, err := MongoDB.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	collection := database.Collection(variable.RevokeTokenColl)
	_, err = collection.DeleteMany(ctx, bson.M{
		"expired_at": bson.M{
			"$lt": currentTime,
		},
	})
	if err != nil {
		log.Println("Error deleting expired tokens:", err)
		return
	}

	// fmt.Println("RevokeTokenExpired execute!", currentTime)
}
