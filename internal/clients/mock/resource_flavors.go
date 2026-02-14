package mock

import (
	"cmyk/internal/models"

	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
)

// ListResourceFlavors reads and parses the mock resource flavors data from JSON file
func (c Client) ListResourceFlavors() ([]models.ResourceFlavor, error) {
	data, err := os.ReadFile("./internal/clients/mock/resource_flavors.json")
	if err != nil {
		return nil, err
	}

	var result []models.ResourceFlavor
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetResourceFlavor reads and returns a mock resource flavor by name
func (c Client) GetResourceFlavor(name string) (*models.ResourceFlavor, error) {
	data, err := os.ReadFile("./internal/clients/mock/resource_flavors.json")
	if err != nil {
		return nil, err
	}

	var flavors []models.ResourceFlavor
	if err := json.Unmarshal(data, &flavors); err != nil {
		return nil, err
	}

	for _, rf := range flavors {
		if rf.Name == name {
			return &rf, nil
		}
	}

	return nil, fiber.ErrNotFound
}
