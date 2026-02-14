package handlers

import (
	"cmyk/internal/models"

	"log"

	"github.com/gofiber/fiber/v2"
)

// ReadPods returns pods as JSON
// @Description Get pods
// @Summary Get pods
// @Tags Pods
// @Produce json
// @Success 200 {array} models.Pod
// @Router /api/v1/pods [get]
func (h Handlers) ReadPods(c *fiber.Ctx) error {
	var pods []models.Pod
	var err error

	if h.EnvClient.IsMockMode() {
		pods, err = h.MockClient.ListPods()
	} else {
		pods, err = h.K8sClient.ListPods()
	}
	if err != nil {
		log.Printf("failed reading pods: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed reading pods"})
	}
	if len(pods) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.JSON(pods)
}

// ReadPodDetail returns pod detail as JSON
// @Description Get pod detail
// @Summary Get pod detail
// @Tags Pods
// @Produce json
// @Param namespace path string true "Pod namespace"
// @Param name path string true "Pod name"
// @Success 200 {object} models.PodDetail
// @Router /api/v1/pods/{namespace}/{name} [get]
func (h Handlers) ReadPodDetail(c *fiber.Ctx) error {
	var podDetail *models.PodDetail
	var err error

	namespace := c.Params("namespace")
	name := c.Params("name")

	if h.EnvClient.IsMockMode() {
		podDetail, err = h.MockClient.GetPod(namespace, name)
	} else {
		podDetail, err = h.K8sClient.GetPod(namespace, name)
	}
	if err != nil {
		log.Printf("failed reading pod: %v", err)
		if err == fiber.ErrNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "Pod not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed reading pod"})
	}
	return c.JSON(podDetail)
}
