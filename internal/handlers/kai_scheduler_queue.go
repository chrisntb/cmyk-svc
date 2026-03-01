package handlers

import (
	"cmyk/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// ReadKaiSchedulerQueues returns kai scheduler parent queues as JSON
// @Description Get kai scheduler parent queues
// @Summary Get kai scheduler parent queues
// @Tags KaiSchedulerQueues
// @Produce json
// @Success 200 {array} models.KaiSchedulerParentQueue
// @Success 204
// @Router /api/v1/kai-scheduler-queues [get]
func (h Handlers) ReadKaiSchedulerQueues(c *fiber.Ctx) error {
	var queues []models.KaiSchedulerParentQueue
	var err error

	if h.EnvClient.IsMockMode() {
		queues, err = h.MockClient.ListKaiSchedulerParentQueues()
	} else {
		queues, err = h.K8sClient.ListKaiSchedulerParentQueues()
	}
	if err != nil {
		log.Printf("failed reading kai scheduler queues: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed reading kai scheduler queues"})
	}
	if len(queues) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.JSON(queues)
}

// ReadKaiSchedulerChildQueues returns child queues for a given parent queue as JSON
// @Description Get child queues for a kai scheduler parent queue
// @Summary Get kai scheduler child queues
// @Tags KaiSchedulerQueues
// @Produce json
// @Param name path string true "Parent queue name"
// @Success 200 {array} models.KaiSchedulerChildQueue
// @Success 204
// @Failure 404 {object} models.Error
// @Router /api/v1/kai-scheduler-queues/{name}/child-queues [get]
func (h Handlers) ReadKaiSchedulerChildQueues(c *fiber.Ctx) error {
	var queues []models.KaiSchedulerChildQueue
	var err error

	name := c.Params("name")

	if h.EnvClient.IsMockMode() {
		queues, err = h.MockClient.GetKaiSchedulerChildQueues(name)
	} else {
		queues, err = h.K8sClient.GetKaiSchedulerChildQueues(name)
	}
	if err != nil {
		log.Printf("failed reading kai scheduler child queues: %v", err)
		if err == fiber.ErrNotFound || apierrors.IsNotFound(err) {
			return c.Status(404).JSON(fiber.Map{"error": "Kai scheduler queue not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed reading kai scheduler child queues"})
	}
	if len(queues) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.JSON(queues)
}
