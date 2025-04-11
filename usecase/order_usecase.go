package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ADV3/repository"
)

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

type Order struct {
	ID          string
	UserID      string
	Items       []OrderItem
	TotalAmount float64
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

type Payment struct {
	OrderID       string
	Amount        float64
	PaymentMethod string
	CardNumber    string
	CardHolder    string
	ExpiryDate    string
	CVV           string
}

type OrderUseCase interface {
	CreateOrder(userID string, items []OrderItem) (*Order, error)
	GetOrder(orderID string) (*Order, error)
	ListOrders(userID string, page, limit int) ([]Order, int, error)
	ProcessPayment(payment Payment) (string, string, error)
}

type orderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewOrderUseCase(orderRepo repository.OrderRepository) OrderUseCase {
	return &orderUseCase{
		orderRepo: orderRepo,
	}
}

func (uc *orderUseCase) CreateOrder(userID string, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	order := &Order{
		ID:          uuid.New().String(),
		UserID:      userID,
		Items:       items,
		TotalAmount: totalAmount,
		Status:      "pending",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	err := uc.orderRepo.SaveOrder(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (uc *orderUseCase) GetOrder(orderID string) (*Order, error) {
	return uc.orderRepo.GetOrderByID(orderID)
}

func (uc *orderUseCase) ListOrders(userID string, page, limit int) ([]Order, int, error) {
	return uc.orderRepo.GetOrdersByUserID(userID, page, limit)
}

func (uc *orderUseCase) ProcessPayment(payment Payment) (string, string, error) {
	// In a real application, this would integrate with a payment processor
	// For this example, we'll just simulate a successful payment
	paymentID := uuid.New().String()
	
	// Update order status to "paid"
	order, err := uc.orderRepo.GetOrderByID(payment.OrderID)
	if err != nil {
		return "", "", err
	}

	order.Status = "paid"
	order.UpdatedAt = time.Now().Format(time.RFC3339)
	
	err = uc.orderRepo.UpdateOrder(order)
	if err != nil {
		return "", "", err
	}

	return paymentID, "success", nil
}
