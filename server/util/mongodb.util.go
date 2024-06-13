package util

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PaginateMongo(model interface{}, collection *mongo.Collection, c *fiber.Ctx, filter bson.M, page, limit int, orderBy, order string) error {
	// Tentukan skip dan limit untuk pagination
	skip := (page - 1) * limit

	// Tentukan opsi pencarian
	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(limit))
	sortOrder := 1
	if order == "desc" {
		sortOrder = -1
	}
	opts.SetSort(bson.D{{Key: orderBy, Value: sortOrder}})

	// Lakukan pencarian
	cursor, err := collection.Find(c.Context(), filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("internal server error %s", err.Error()),
		})
	}
	defer cursor.Close(c.Context())

	// Dekode hasil pencarian ke dalam slice dari map
	var results []bson.M
	if err := cursor.All(c.Context(), &results); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("error decoding results %s", err.Error()),
		})
	}

	// Hapus field password dari setiap object
	for _, result := range results {
		delete(result, "password")
	}

	// Inisialisasi results sebagai slice kosong jika tidak ada data
	if len(results) == 0 {
		results = []bson.M{}
	}

	// Hitung total dokumen
	total, err := collection.CountDocuments(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("error counting documents %s", err.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":  results,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}
