package handlers

import (
	"cmyk/internal/models"

	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// ReadResourceFlavors returns resource flavors as JSON
// @Description Get resource flavors
// @Summary Get resource flavors
// @Tags ResourceFlavors
// @Produce json
// @Success 200 {array} models.ResourceFlavor
// @Router /api/v1/resource-flavors [get]
func (h Handlers) ReadResourceFlavors(c *fiber.Ctx) error {
	var flavors []models.ResourceFlavor
	var err error

	if h.EnvClient.IsMockMode() {
		flavors, err = h.MockClient.ListResourceFlavors()
	} else {
		flavors, err = h.K8sClient.ListResourceFlavors()
	}
	if err != nil {
		log.Printf("failed reading resource flavors: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed reading resource flavors"})
	}
	if len(flavors) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.JSON(flavors)
}

// ReadResourceFlavorDetail returns resource flavor detail as JSON
// @Description Get resource flavor detail
// @Summary Get resource flavor detail
// @Tags ResourceFlavors
// @Produce json
// @Param name path string true "ResourceFlavor name"
// @Success 200 {object} models.ResourceFlavor
// @Failure 404 {object} models.Error
// @Router /api/v1/resource-flavors/{name} [get]
func (h Handlers) ReadResourceFlavorDetail(c *fiber.Ctx) error {
	var flavor *models.ResourceFlavor
	var err error

	name := c.Params("name")

	if h.EnvClient.IsMockMode() {
		flavor, err = h.MockClient.GetResourceFlavor(name)
	} else {
		flavor, err = h.K8sClient.GetResourceFlavor(name)
	}
	if err != nil {
		log.Printf("failed reading resource flavor: %v", err)
		if err == fiber.ErrNotFound || apierrors.IsNotFound(err) {
			return c.Status(404).JSON(fiber.Map{"error": "Resource flavor not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed reading resource flavor"})
	}
	return c.JSON(flavor)
}

// CreateResourceFlavor creates a new resource flavor
// @Description Create a new resource flavor
// @Summary Create resource flavor
// @Tags ResourceFlavors
// @Accept json
// @Produce json
// @Param resourceFlavor body models.ResourceFlavor true "ResourceFlavor to create"
// @Success 201 {object} models.ResourceFlavor
// @Failure 400 {object} models.Error
// @Router /api/v1/resource-flavors [post]
func (h Handlers) CreateResourceFlavor(c *fiber.Ctx) error {
	var rawBody map[string]any
	if err := c.BodyParser(&rawBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Code:    fiber.StatusBadRequest,
			Message: utils.StatusMessage(fiber.StatusBadRequest),
			Reason:  "Cannot parse JSON",
		})
	}

	rf, err := validateResourceFlavorSchema(rawBody)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Code:    fiber.StatusBadRequest,
			Message: utils.StatusMessage(fiber.StatusBadRequest),
			Reason:  err.Error(),
		})
	}

	if h.K8sClient != nil {
		created, err := h.K8sClient.CreateResourceFlavor(*rf)
		if err != nil {
			log.Printf("failed creating resource flavor: %v", err)
			return c.Status(fiber.StatusBadGateway).JSON(models.Error{
				Code:    fiber.StatusBadGateway,
				Message: utils.StatusMessage(fiber.StatusBadGateway),
				Reason:  err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(created)
	}

	return c.Status(fiber.StatusCreated).JSON(rf)
}

// DeleteResourceFlavor deletes a resource flavor
// @Description Delete a resource flavor
// @Summary Delete resource flavor
// @Tags ResourceFlavors
// @Produce json
// @Param name path string true "ResourceFlavor name"
// @Success 204
// @Failure 502 {object} models.Error
// @Router /api/v1/resource-flavors/{name} [delete]
func (h Handlers) DeleteResourceFlavor(c *fiber.Ctx) error {
	name := c.Params("name")

	if h.K8sClient != nil {
		if err := h.K8sClient.DeleteResourceFlavor(name); err != nil {
			log.Printf("failed deleting resource flavor: %v", err)
			return c.Status(fiber.StatusBadGateway).JSON(models.Error{
				Code:    fiber.StatusBadGateway,
				Message: utils.StatusMessage(fiber.StatusBadGateway),
				Reason:  err.Error(),
			})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

//revive:disable:cyclomatic
func validateResourceFlavorSchema(rawBody map[string]any) (*models.ResourceFlavor, error) {
	nameValue, exists := rawBody["name"]
	if !exists {
		return nil, fmt.Errorf("field 'name' is required")
	}

	name, ok := nameValue.(string)
	if !ok {
		return nil, fmt.Errorf("field 'name' must be a string")
	}

	if name == "" {
		return nil, fmt.Errorf("field 'name' cannot be empty")
	}

	rf := &models.ResourceFlavor{Name: name}

	if nodeLabels, exists := rawBody["nodeLabels"]; exists {
		labels, ok := nodeLabels.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("field 'nodeLabels' must be an object")
		}
		rf.NodeLabels = make(map[string]string, len(labels))
		for k, v := range labels {
			s, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("nodeLabels values must be strings")
			}
			rf.NodeLabels[k] = s
		}
	}

	if nodeTaints, exists := rawBody["nodeTaints"]; exists {
		taints, ok := nodeTaints.([]any)
		if !ok {
			return nil, fmt.Errorf("field 'nodeTaints' must be an array")
		}
		for _, item := range taints {
			t, ok := item.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("nodeTaints items must be objects")
			}
			taint := models.NodeTaint{}
			if v, ok := t["key"].(string); ok {
				taint.Key = v
			}
			if v, ok := t["value"].(string); ok {
				taint.Value = v
			}
			if v, ok := t["effect"].(string); ok {
				taint.Effect = v
			}
			rf.NodeTaints = append(rf.NodeTaints, taint)
		}
	}

	if tolerations, exists := rawBody["tolerations"]; exists {
		tols, ok := tolerations.([]any)
		if !ok {
			return nil, fmt.Errorf("field 'tolerations' must be an array")
		}
		for _, item := range tols {
			t, ok := item.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("tolerations items must be objects")
			}
			tol := models.Toleration{}
			if v, ok := t["key"].(string); ok {
				tol.Key = v
			}
			if v, ok := t["operator"].(string); ok {
				tol.Operator = v
			}
			if v, ok := t["value"].(string); ok {
				tol.Value = v
			}
			if v, ok := t["effect"].(string); ok {
				tol.Effect = v
			}
			rf.Tolerations = append(rf.Tolerations, tol)
		}
	}

	if topologyName, exists := rawBody["topologyName"]; exists {
		tn, ok := topologyName.(string)
		if !ok {
			return nil, fmt.Errorf("field 'topologyName' must be a string")
		}
		rf.TopologyName = tn
	}

	return rf, nil
}

//revive:enable:cyclomatic
