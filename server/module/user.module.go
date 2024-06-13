package module

import (
	"backend/server/connection"
	"backend/server/env"
	"backend/server/middleware"
	"backend/server/model"
	"backend/server/util"
	"backend/server/variable"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{}

func (ref User) Route(api fiber.Router) {
	handler := UserHandler{}
	route := api.Group("/user")

	route.Get("/detail", middleware.UseAuth, middleware.RoleAccess([]string{
		variable.AdministratorRole,
		variable.PartaiRole,
		variable.PelaksanaRole,
		variable.SaksiRole,
		variable.KandidatRole,
	}), handler.Detail)

	// manage
	route.Get("/paginate", middleware.UseAuth, middleware.RoleAccess([]string{
		variable.AdministratorRole,
		variable.PartaiRole,
		variable.PelaksanaRole,
	}), handler.Paginate)
	route.Post("/new", middleware.UseAuth, middleware.RoleAccess([]string{
		variable.AdministratorRole,
		variable.PartaiRole,
		variable.PelaksanaRole,
	}), handler.New)
	route.Put("/edit/:id", middleware.UseAuth, middleware.RoleAccess([]string{
		variable.AdministratorRole,
		variable.PartaiRole,
		variable.PelaksanaRole,
	}), handler.Edit)
	route.Patch("/change/:id", middleware.UseAuth, middleware.RoleAccess([]string{
		variable.AdministratorRole,
		variable.PartaiRole,
		variable.PelaksanaRole,
	}), handler.Change)
	route.Delete("/remove/:id", middleware.UseAuth, middleware.RoleAccess([]string{
		variable.AdministratorRole,
		variable.PartaiRole,
		variable.PelaksanaRole,
	}), handler.Remove)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type UserHandler struct{}

func (handler UserHandler) Detail(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":         user.ID.Hex(),
		"name":       user.Name,
		"image_url":  user.ImageURL,
		"username":   user.Username,
		"role_name":  user.RoleID,
		"created_at": user.CreatedAt,
	})
}

func (handler UserHandler) Paginate(c *fiber.Ctx) error {
	// Ambil query parameter
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")
	search := c.Query("search", "")
	orderBy := c.Query("order_by", "created_at")
	order := c.Query("order", "asc")

	// Konversi page dan limit ke integer
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Panggil fungsi Paginate dengan parameter-parameter yang diperlukan
	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("connect mongodb %s", err.Error()),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"username": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	collection := database.Collection(variable.UserColl)
	return util.PaginateMongo(model.User{}, collection, c, filter, page, limit, orderBy, order)
}

func (handler UserHandler) New(c *fiber.Ctx) error {
	claims := c.Locals("claims").(jwt.MapClaims)
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user_id is not valid",
		})
	}

	// Ambil data pengguna dari request body
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}
	if user.Name == "" || user.Username == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "missing required fields",
		})
	}

	// Inisialisasi MongoDB
	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("connect mongodb %s", err.Error()),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())

	// Buat objek pengguna baru
	user.ID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := database.Collection(variable.UserColl)
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to insert user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "user created successfully",
		"user_id":     user.ID.Hex(),
		"inserted_by": userIDStr,
	})
}

func (handler UserHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals("claims").(jwt.MapClaims)
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user_id is not valid",
		})
	}

	// Ambil id pengguna yang akan diedit dari URL parameter
	id := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user ID",
		})
	}

	// Ambil data yang akan diupdate dari request body
	var updateData model.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	// Inisialisasi MongoDB
	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("connect mongodb %s", err.Error()),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	collection := database.Collection(variable.UserColl)

	// Buat dokumen pembaruan
	update := bson.M{
		"$set": bson.M{
			"name":       updateData.Name,
			"image_url":  updateData.ImageURL,
			"email":      updateData.Email,
			"username":   updateData.Username,
			"is_verify":  updateData.IsVerify,
			"role_id":    updateData.RoleID,
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": userID,
	}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "user updated successfully",
		"updated_by": userIDStr,
	})
}

func (handler UserHandler) Change(c *fiber.Ctx) error {
	claims := c.Locals("claims").(jwt.MapClaims)
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user_id is not valid",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"changed_by": userIDStr,
	})
}

func (handler UserHandler) Remove(c *fiber.Ctx) error {
	claims := c.Locals("claims").(jwt.MapClaims)
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user_id is not valid",
		})
	}

	// Ambil id pengguna yang akan dihapus dari URL parameter
	id := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user ID",
		})
	}

	// Inisialisasi MongoDB
	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("connect mongodb %s", err.Error()),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	collection := database.Collection(variable.UserColl)

	// Lakukan penghapusan
	_, err = collection.DeleteOne(ctx, bson.M{
		"_id": userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "user deleted successfully",
		"deleted_by": userIDStr,
	})
}
