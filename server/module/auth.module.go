package module

import (
	"backend/server/connection"
	"backend/server/env"
	"backend/server/middleware"
	"backend/server/model"
	"backend/server/util"
	"backend/server/variable"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct{}

func (ref Auth) Route(api fiber.Router) {
	handler := AuthHandler{}
	route := api.Group("/auth")

	// route.Post("/register", handler.Register)
	route.Post("/forgot-password", handler.ForgotPassword)

	// JWT Handler with JTI
	route.Post("/login", handler.Login)
	route.Delete("/logout", middleware.UseAuth, handler.Logout)
	route.Get("/token-validation", middleware.UseAuth, handler.TokenValidation)
}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type AuthHandler struct{}

func (handler AuthHandler) Register(c *fiber.Ctx) error {
	var body model.UserRegisterBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("connect mongodb %s", err.Error()),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	user := model.User{}
	collection := database.Collection(variable.UserColl)

	err = collection.FindOne(ctx, bson.M{
		"username":  body.Username,
		"is_verify": true,
	}).Decode(&user)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "username already exists",
		})
	} else if err != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error on query",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not hash password",
		})
	}

	NowAt := primitive.NewDateTimeFromTime(time.Now())

	err = collection.FindOne(ctx, bson.M{
		"username": body.Username,
	}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = collection.InsertOne(ctx, bson.M{
				"name": body.Name,

				"username":  body.Username,
				"password":  string(hashedPassword),
				"is_verify": false,

				"created_at": NowAt,
			})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "cannot inserted",
				})
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error on query",
			})
		}
	} else {
		_, err = collection.UpdateOne(ctx, bson.M{
			"username": body.Username,
		}, bson.M{
			"$set": model.User{
				Name:      body.Name,
				Password:  string(hashedPassword),
				UpdatedAt: &NowAt,
			},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "cannot update",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user pending",
	})
}

func (handler AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var err error

	var body model.UserForgotPasswordBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("connect mongodb %s", err.Error()),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	var otp_ref string
	var otp_code string
	user := model.User{}
	collection := database.Collection(variable.UserColl)
	if body.Email != nil {
		err = collection.FindOne(ctx, bson.M{
			"email": body.Email,
		}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "account not found",
				})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "internal server error",
				})
			}
		}
	} else if body.Username != nil {
		err = collection.FindOne(ctx, bson.M{
			"username": *body.Username,
		}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "account not found",
				})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "internal server error",
				})
			}
		}
	} else if body.Email != nil && body.Username != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "cannot use email and username same time",
		})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	// In a real application, you would send the OTP via email
	// For demonstration, we just return it in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "OTP sent",
		"otp_ref":  otp_ref,
		"otp_code": otp_code,
	})
}

func (handler AuthHandler) Login(c *fiber.Ctx) error {
	browser_id := c.Locals("browser_id").(string)

	var body model.UserLoginBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	// headers := c.Request().Header.RawHeaders()
	// fmt.Println("headers:", string(headers))

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("connect mongodb %s", err.Error()),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	message_not_found := "username or password is wrong"

	user := model.User{}
	collection := database.Collection(variable.UserColl)
	err = collection.FindOne(ctx, bson.M{
		"username": body.Username,
	}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": message_not_found,
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
	}

	Encryption := util.Encryption{}
	decodePassword, err := Encryption.DecodeWithSecret(user.Password)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error decryption",
		})
	}
	if body.Password != decodePassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": message_not_found,
		})
	}

	user_id := user.ID.Hex()
	statusCode, err := checkJti(database, ctx, bson.M{
		"user_id": user_id,
	})
	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	token, statusCode, err := generateToken(database, ctx, c.Get("user-agent"), user_id, browser_id)
	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func (handler AuthHandler) Logout(c *fiber.Ctx) error {
	browser_id := c.Locals("browser_id").(string)
	claims := c.Locals("claims").(jwt.MapClaims)
	user_id := claims["user_id"].(string)
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
	defer database.Client().Disconnect(ctx)

	collection := database.Collection(variable.RevokeTokenColl)
	collection.DeleteOne(ctx, bson.M{
		"user_id":    user_id,
		"jwt_id":     jti,
		"browser_id": browser_id,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success logout",
	})
}

func (handler AuthHandler) TokenValidation(c *fiber.Ctx) error {
	claims := c.Locals("claims").(jwt.MapClaims)
	userID, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user_id is not valid",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id": userID,
	})
}

// ----------------------------------------------------------------

func insertLogin(database *mongo.Database, ctx context.Context, user_agent string, user_id string, browser_id string, jti string) (int, error) {
	var err error
	var collection *mongo.Collection

	setting := model.Setting{}
	expired_login := variable.KeyExpiredLoginJwt
	collection = database.Collection(variable.SettingColl)
	err = collection.FindOne(ctx, bson.M{
		"key": expired_login,
	}).Decode(&setting)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.StatusInternalServerError, fmt.Errorf("setting not found")
		} else {
			return fiber.StatusInternalServerError, fmt.Errorf("internal server error")
		}
	}

	expiredLoginInt, err := strconv.ParseInt(setting.Value, 10, 64)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf(fmt.Sprintf("Error parsing int64: %s", err.Error()))
	}
	duration := time.Duration(expiredLoginInt)
	expired_login_time := time.Now().Add(duration)
	ExpiredAt := primitive.NewDateTimeFromTime(expired_login_time) // microsecond
	NowAt := primitive.NewDateTimeFromTime(time.Now())

	collection = database.Collection(variable.RevokeTokenColl)
	_, err = collection.InsertOne(ctx, bson.M{
		"user_id":    user_id,
		"jwt_id":     jti,
		"browser_id": browser_id,
		"expired_at": ExpiredAt,
		"login_at":   NowAt,
	})
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("cannot inserted revoke_token")
	}

	collection = database.Collection(variable.LoginHistoryColl)
	_, err = collection.InsertOne(ctx, bson.M{
		"user_id":    user_id,
		"user_agent": user_agent,
		"login_at":   NowAt,
	})
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("cannot inserted login_history")
	}

	return 0, nil
}

func generateToken(database *mongo.Database, ctx context.Context, user_agent string, user_id string, browser_id string) (string, int, error) {
	var err error

	JWT := util.JWT{}
	token, jti, err := JWT.Generate(user_id) // Short-lived token for access
	if err != nil {
		return "", fiber.StatusInternalServerError, fmt.Errorf("could not generate token")
	}

	statusCode, err := insertLogin(database, ctx, user_agent, user_id, browser_id, jti)
	if err != nil {
		return "", statusCode, err
	}
	return token, 0, nil
}

func checkJti(database *mongo.Database, ctx context.Context, filter bson.M) (int, error) {
	var err error
	var collection *mongo.Collection

	collection = database.Collection(variable.RevokeTokenColl)
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("internal server error")
	}

	setting := model.Setting{}
	max_login := variable.KeyMaxLoginJwt
	collection = database.Collection(variable.SettingColl)
	err = collection.FindOne(ctx, bson.M{
		"key": max_login,
	}).Decode(&setting)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.StatusInternalServerError, fmt.Errorf("setting not found")
		} else {
			return fiber.StatusInternalServerError, fmt.Errorf("internal server error")
		}
	}

	maxAttempts, err := strconv.Atoi(setting.Value)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("error parsing maximum attempts value: %v", err)
	}
	countInt := int(count)
	if countInt >= maxAttempts {
		return fiber.StatusForbidden, fmt.Errorf("maximum login attempts exceeded")
	}

	return 0, nil
}
