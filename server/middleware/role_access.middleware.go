package middleware

import (
	"backend/server/model"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func RoleAccess(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(model.Role)
		userRole := role.Name
		allowed := false
		for _, role := range allowedRoles {
			if role == userRole {
				allowed = true
				break
			}
		}
		if allowed {
			c.Locals("role_name", userRole)
			return c.Next()
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": fmt.Sprintf("Role %s cannot access", userRole),
		})
	}
}
