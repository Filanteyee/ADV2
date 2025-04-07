package handler

import (
	"net/http"
	"order_service/repository"
	"order_service/usecase"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	usecase *usecase.OrderUsecase
}

func NewOrderHandler(uc *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: uc}
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.usecase.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении заказов"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.usecase.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заказ не найден"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}
	if err := h.usecase.UpdateStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении статуса"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении заказа"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req repository.Order
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	id, err := h.usecase.CreateOrder(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании заказа"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Order created",
		"order_id": id,
	})
}
