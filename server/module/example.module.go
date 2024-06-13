package module

import (
	"backend/server/util"

	"github.com/gofiber/fiber/v2"
)

type Example struct{}

func (ref Example) Route(api fiber.Router) {
	handler := ExampleHandler{}
	route := api.Group("/example")

	route.Get("/", handler.HelloWorld)
	route.Get("/trigger/:value", handler.Trigger)
	route.Get("/encode/:browser_id", handler.Encode)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type ExampleHandler struct{}

func (handler ExampleHandler) HelloWorld(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Hello World",
	})
}

func (handler ExampleHandler) Trigger(c *fiber.Ctx) error {
	// var err error

	subdomain, _ := c.Locals("subdomain").(string)
	browserID, _ := c.Locals("browser_id").(string)
	partaiID, _ := c.Locals("partai_id").(string)

	value := c.Params("value")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "OK",
		"value":      value,
		"subdomain":  subdomain,
		"browser_id": browserID,
		"partai_id":  partaiID,
	})
}

func (handler ExampleHandler) Encode(c *fiber.Ctx) error {
	var err error

	host := c.Hostname()
	browser_id := c.Params("browser_id")

	Encryption := util.Encryption{}
	browser_id, err = Encryption.Encode(host, browser_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "OK",
		"browser_id": browser_id,
	})
}
