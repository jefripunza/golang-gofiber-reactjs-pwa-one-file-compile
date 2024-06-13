package middleware

import (
	"backend/server/connection"
	"backend/server/env"
	"backend/server/model"
	"backend/server/variable"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UseAuth(c *fiber.Ctx) error {
	secretKey := env.GetSecretKey()
	browser_id := c.Locals("browser_id").(string)
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header missing",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Authorization header format",
		})
	}

	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "user_id is not valid",
			})
		}
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "user_id is not a valid ObjectID",
			})
		}
		jti := claims["jti"].(string)

		MongoDB := connection.MongoDB{}
		client, ctx, err := MongoDB.Connect()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("connect mongodb %s", err.Error()),
			})
		}
		defer client.Disconnect(ctx)
		database := client.Database(env.GetMongoName())
		var collection *mongo.Collection

		exist := model.RevokeToken{}
		collection = database.Collection(variable.RevokeTokenColl)
		err = collection.FindOne(ctx, bson.M{
			"user_id":    userIDStr,
			"jwt_id":     jti,
			"browser_id": browser_id,
		}).Decode(&exist)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "token cannot be used",
				})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "internal server error",
				})
			}
		}

		var user model.User
		collection = database.Collection(variable.UserColl)
		err = collection.FindOne(ctx, bson.M{
			"_id": userID,
		}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "user not found",
				})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "internal server error",
				})
			}
		}

		roleID, err := primitive.ObjectIDFromHex(user.RoleID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "role_id is not a valid ObjectID",
			})
		}
		var role model.Role
		collection = database.Collection(variable.RoleColl)
		err = collection.FindOne(ctx, bson.M{
			"_id": roleID,
		}).Decode(&role)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "role not found",
				})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "internal server error",
				})
			}
		}

		// Add claims to context
		c.Locals("claims", claims)
		c.Locals("user", user)
		c.Locals("role", role)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	return c.Next()
}
