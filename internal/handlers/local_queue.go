package handlers

import (
	"cmyk/internal/models"

	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// ReadLocalQueues returns local queues as JSON
// @Description Get local queues
// @Summary Get local queues
// @Tags LocalQueues
// @Produce json
// @Success 200 {array} models.LocalQueue
// @Router /api/v1/local-queues [get]
func (h Handlers) ReadLocalQueues(c *fiber.Ctx) error {
	var queues []models.LocalQueue
	var err error

	if h.EnvClient.IsMockMode() {
		queues, err = h.MockClient.ListLocalQueues()
	} else {
		queues, err = h.K8sClient.ListLocalQueues()
	}
	if err != nil {
		log.Printf("failed reading local queues: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed reading local queues"})
	}
	if len(queues) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.JSON(queues)
}

// ReadLocalQueueDetail returns local queue detail as JSON
// @Description Get local queue detail
// @Summary Get local queue detail
// @Tags LocalQueues
// @Produce json
// @Param namespace path string true "LocalQueue namespace"
// @Param name path string true "LocalQueue name"
// @Success 200 {object} models.LocalQueue
// @Failure 404 {object} models.Error
// @Router /api/v1/namespaces/{namespace}/local-queues/{name} [get]
func (h Handlers) ReadLocalQueueDetail(c *fiber.Ctx) error {
	var queue *models.LocalQueue
	var err error

	namespace := c.Params("namespace")
	name := c.Params("name")

	if h.EnvClient.IsMockMode() {
		queue, err = h.MockClient.GetLocalQueue(namespace, name)
	} else {
		queue, err = h.K8sClient.GetLocalQueue(namespace, name)
	}
	if err != nil {
		log.Printf("failed reading local queue: %v", err)
		if err == fiber.ErrNotFound || apierrors.IsNotFound(err) {
			return c.Status(404).JSON(fiber.Map{"error": "Local queue not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed reading local queue"})
	}
	return c.JSON(queue)
}

// CreateLocalQueue creates a new local queue
// @Description Create a new local queue
// @Summary Create local queue
// @Tags LocalQueues
// @Accept json
// @Produce json
// @Param namespace path string true "LocalQueue namespace"
// @Param localQueue body models.LocalQueue true "LocalQueue to create"
// @Success 201 {object} models.LocalQueue
// @Failure 400 {object} models.Error
// @Router /api/v1/namespaces/{namespace}/local-queues [post]
func (h Handlers) CreateLocalQueue(c *fiber.Ctx) error {
	namespace := c.Params("namespace")

	var rawBody map[string]any
	if err := c.BodyParser(&rawBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Code:    fiber.StatusBadRequest,
			Message: utils.StatusMessage(fiber.StatusBadRequest),
			Reason:  "Cannot parse JSON",
		})
	}

	lq, err := validateLocalQueueSchema(rawBody)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Code:    fiber.StatusBadRequest,
			Message: utils.StatusMessage(fiber.StatusBadRequest),
			Reason:  err.Error(),
		})
	}

	if h.K8sClient != nil {
		created, err := h.K8sClient.CreateLocalQueue(namespace, *lq)
		if err != nil {
			log.Printf("failed creating local queue: %v", err)
			return c.Status(fiber.StatusBadGateway).JSON(models.Error{
				Code:    fiber.StatusBadGateway,
				Message: utils.StatusMessage(fiber.StatusBadGateway),
				Reason:  err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(created)
	}

	lq.Namespace = namespace
	return c.Status(fiber.StatusCreated).JSON(lq)
}

// DeleteLocalQueue deletes a local queue
// @Description Delete a local queue
// @Summary Delete local queue
// @Tags LocalQueues
// @Produce json
// @Param namespace path string true "LocalQueue namespace"
// @Param name path string true "LocalQueue name"
// @Success 204
// @Failure 502 {object} models.Error
// @Router /api/v1/namespaces/{namespace}/local-queues/{name} [delete]
func (h Handlers) DeleteLocalQueue(c *fiber.Ctx) error {
	namespace := c.Params("namespace")
	name := c.Params("name")

	if h.K8sClient != nil {
		if err := h.K8sClient.DeleteLocalQueue(namespace, name); err != nil {
			log.Printf("failed deleting local queue: %v", err)
			return c.Status(fiber.StatusBadGateway).JSON(models.Error{
				Code:    fiber.StatusBadGateway,
				Message: utils.StatusMessage(fiber.StatusBadGateway),
				Reason:  err.Error(),
			})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func validateLocalQueueSchema(rawBody map[string]any) (*models.LocalQueue, error) {
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

	cqValue, exists := rawBody["clusterQueue"]
	if !exists {
		return nil, fmt.Errorf("field 'clusterQueue' is required")
	}

	clusterQueue, ok := cqValue.(string)
	if !ok {
		return nil, fmt.Errorf("field 'clusterQueue' must be a string")
	}

	if clusterQueue == "" {
		return nil, fmt.Errorf("field 'clusterQueue' cannot be empty")
	}

	lq := &models.LocalQueue{
		Name:         name,
		ClusterQueue: clusterQueue,
	}

	if stopPolicy, exists := rawBody["stopPolicy"]; exists {
		sp, ok := stopPolicy.(string)
		if !ok {
			return nil, fmt.Errorf("field 'stopPolicy' must be a string")
		}
		lq.StopPolicy = sp
	}

	return lq, nil
}
