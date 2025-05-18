package handlers

import (
	"sika/service"

	"github.com/gofiber/fiber/v2"
)

func GetUserByID(userService *service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("UserID")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "userID is required",
			})
		}

		user, err := userService.GetUserByID(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "user not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(user)
	}
}
