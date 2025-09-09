package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"presign-service/config"
	"presign-service/handlers"
	"presign-service/middleware"
	"presign-service/models"
	"presign-service/service"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	s3Service, err := service.NewS3Service(service.S3Config{
		BucketName:  cfg.BucketName,
		Region:      cfg.AWSRegion,
		AccessKeyID: cfg.AWSAccessKeyID,
		SecretKey:   cfg.AWSSecretKey,
		EndpointURL: cfg.AWSEndpointURL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize S3 service: %v", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.Printf("Error: %v", err)
			return c.Status(code).JSON(models.ErrorResponse{
				Error: err.Error(),
			})
		},
	})

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	app.Use(recover.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(models.HealthResponse{
			Status:  "healthy",
			Service: "presign-service",
		})
	})

	signHandler := handlers.NewSignHandler(s3Service)

	authMiddleware := middleware.NewAuthMiddleware(cfg.APIKey)

	app.Post("/sign", authMiddleware, signHandler.Handle)

	log.Printf("Starting presign-service on port %s", cfg.Port)
	log.Printf("Bucket: %s", cfg.BucketName)
	log.Printf("Health check available at: http://localhost:%s/health", cfg.Port)

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
