package main

import (
	"log"
	"os"

	"github.com/Filanteyee/ADV3/repository"
	"github.com/Filanteyee/ADV3/server"
	"github.com/Filanteyee/ADV3/usecase"
)

func main() {
	// Initialize repository
	orderRepo := repository.NewOrderRepository()

	// Initialize use case
	orderUseCase := usecase.NewOrderUseCase(orderRepo)

	// Start gRPC server
	port := os.Getenv("ORDER_SERVICE_PORT")
	if port == "" {
		port = "50051" // Default port
	}

	if err := server.StartOrderServer(orderUseCase, port); err != nil {
		log.Fatalf("Failed to start Order Service: %v", err)
	}
}
