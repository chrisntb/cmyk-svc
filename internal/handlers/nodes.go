package handlers

import (
	"cmyk/internal/models"

	"log"

	"github.com/gofiber/fiber/v2"
)

// ReadNodes returns nodes as JSON
// @Description Get nodes
// @Summary Get nodes
// @Tags Nodes
// @Produce json
// @Success 200 {array} models.Node
// @Router /api/v1/nodes [get]
func (h Handlers) ReadNodes(c *fiber.Ctx) error {
	var nodes []models.Node
	var err error

	if h.EnvClient.IsMockMode() {
		nodes, err = h.MockClient.ListNodes()
	} else {
		nodes, err = h.K8sClient.ListNodes()
	}
	if err != nil {
		log.Printf("failed reading nodes: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed reading nodes"})
	}
	if len(nodes) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.JSON(nodes)
}

// ReadNodeDetail returns node detail as JSON
// @Description Get node detail
// @Summary Get node detail
// @Tags Nodes
// @Produce json
// @Param name path string true "Node name"
// @Success 200 {object} models.NodeDetail
// @Router /api/v1/nodes/{name} [get]
func (h Handlers) ReadNodeDetail(c *fiber.Ctx) error {
	var nodeDetail *models.NodeDetail
	var err error

	name := c.Params("name")

	if h.EnvClient.IsMockMode() {
		nodeDetail, err = h.MockClient.GetNode(name)
	} else {
		nodeDetail, err = h.K8sClient.GetNode(name)
	}
	if err != nil {
		log.Printf("failed reading node: %v", err)
		if err == fiber.ErrNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "Node not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed reading node"})
	}
	return c.JSON(nodeDetail)
}
