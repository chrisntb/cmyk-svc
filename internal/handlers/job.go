package handlers

import (
	"cmyk/internal/models"

	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateJob func creates a new job
// @Description Create a new job
// @Summary Create job
// @Tags Job
// @Accept json
// @Produce json
// @Param job body models.Job true "Job to create"
// @Success 201 {object} models.Job
// @Failure 400 {object} models.Error
// @Router /api/v1/jobs [post]
func (h Handlers) CreateJob(c *fiber.Ctx) error {
	var rawBody map[string]any
	if err := c.BodyParser(&rawBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Code:    fiber.StatusBadRequest,
			Message: utils.StatusMessage(fiber.StatusBadRequest),
			Reason:  "Cannot parse JSON",
		})
	}

	job, err := validateJobSchema(rawBody)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Code:    fiber.StatusBadRequest,
			Message: utils.StatusMessage(fiber.StatusBadRequest),
			Reason:  err.Error(),
		})
	}

	// TODO - Remove this check when the service is changed to not start if the K8s client cannot be initialized properly
	if h.K8sClient != nil {
		err = h.runJobInCluster(job.Name)
		if err != nil {
			log.Printf("failed running job in cluster: %v", err)
			return c.Status(fiber.StatusBadGateway).JSON(models.Error{
				Code:    fiber.StatusBadGateway,
				Message: utils.StatusMessage(fiber.StatusBadGateway),
				Reason:  err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(job)
}

// validateJobSchema validates the job creation request
func validateJobSchema(rawBody map[string]any) (*models.Job, error) {
	if len(rawBody) != 1 {
		return nil, fmt.Errorf("request must contain only 'name' field")
	}

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

	return &models.Job{Name: name}, nil
}

func (h Handlers) runJobInCluster(jobName string) error {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"app": "job",
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:  "job-container",
					Image: "busybox:latest",
					Command: []string{
						"sh",
						"-c",
						"echo 'Job started' && sleep 10 && echo 'Job completed'",
					},
				},
			},
		},
	}

	_, err := h.K8sClient.Clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create pod: %w", err)
	}

	return nil
}
