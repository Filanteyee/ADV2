// main.go
package main

import (
	"context"
	"log"
	"order_service/handler"
	"order_service/repository"
	"order_service/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	db := client.Database("orders_db")

	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/static/orders.html")
	})

	repo := repository.NewOrderRepository(db)
	uc := usecase.NewOrderUsecase(repo)
	h := handler.NewOrderHandler(uc)

	r.POST("/orders", h.CreateOrder)
	r.GET("/orders", h.GetOrders)
	r.GET("/orders/:id", h.GetOrderByID)
	r.PATCH("/orders/:id", h.UpdateOrderStatus)
	r.DELETE("/orders/:id", h.DeleteOrder)

	r.Run(":8080")
}
