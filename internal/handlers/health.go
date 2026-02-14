package handlers

import (
	"cmyk/internal/models"

	"github.com/gofiber/fiber/v2"
)

// Health func returns service ehealth
// @Description Read health
// @Summary Read health
// @Tags Health
// @Produce json
// @Success 200 {object} models.Health
// @Router / [get]
func (h Handlers) Health(c *fiber.Ctx) error {
	return c.JSON(models.Health{Status: "UP"})
}
