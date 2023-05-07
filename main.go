package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/midtrans/midtrans-go"
	"log"
	mg "payment/midtrans_gateway"
)

func main() {
	midtransGateway := mg.NewPayment(&mg.Config{
		ServerKey: "YOUR_SERVER_KEY",
		Env:       midtrans.Sandbox,
	})

	s := fiber.New()

	s.Use(cors.New())
	s.Use(recover.New())
	s.Use(logger.New())

	s.Post("/api/transactions", func(c *fiber.Ctx) error {
		var body mg.CustomerDetailsRequest
		if err := c.BodyParser(&body); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, err)
		}

		response, err := midtransGateway.CreateTransaction(&body)
		if err != nil {
			return errorResponse(c, fiber.StatusInternalServerError, err)
		}

		return successResponse(c, fiber.StatusOK, response)
	})

	s.Get("/api/transactions/status/:orderID", func(c *fiber.Ctx) error {
		orderID := c.Params("orderID")

		transactionStatus, err := midtransGateway.TransactionStatus(orderID)
		if err != nil {
			return errorResponse(c, fiber.StatusInternalServerError, err)
		}

		return successResponse(c, fiber.StatusOK, transactionStatus)
	})

	log.Println("Server is running on port 3000")
	err := s.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func errorResponse(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func successResponse(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"data": data,
	})
}
