package mock

import (
	"cmyk/internal/models"

	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
)

// ListLocalQueues reads and parses the mock local queues data from JSON file
func (c Client) ListLocalQueues() ([]models.LocalQueue, error) {
	data, err := os.ReadFile("./internal/clients/mock/local_queues.json")
	if err != nil {
		return nil, err
	}

	var result []models.LocalQueue
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetLocalQueue reads and returns a mock local queue by namespace and name
func (c Client) GetLocalQueue(namespace, name string) (*models.LocalQueue, error) {
	data, err := os.ReadFile("./internal/clients/mock/local_queues.json")
	if err != nil {
		return nil, err
	}

	var queues []models.LocalQueue
	if err := json.Unmarshal(data, &queues); err != nil {
		return nil, err
	}

	for _, lq := range queues {
		if lq.Namespace == namespace && lq.Name == name {
			return &lq, nil
		}
	}

	return nil, fiber.ErrNotFound
}
