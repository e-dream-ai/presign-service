package handlers

import (
	"context"
	"log"

	"presign-service/models"
	"presign-service/service"

	"github.com/gofiber/fiber/v2"
)

type SignHandler struct {
	storageService *service.StorageService
}

func NewSignHandler(storageService *service.StorageService) *SignHandler {
	return &SignHandler{
		storageService: storageService,
	}
}

func (h *SignHandler) Handle(c *fiber.Ctx) error {
	var req models.SignRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid JSON format",
		})
	}

	if len(req.Keys) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Keys array cannot be empty",
		})
	}

	ctx := context.Background()
	urls, err := h.storageService.PresignMultipleObjects(ctx, req.Keys)
	if err != nil {
		log.Printf("Failed to presign URLs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to generate presigned URLs",
		})
	}

	response := models.SignResponse{
		URLs: urls,
	}

	return c.JSON(response)
}
