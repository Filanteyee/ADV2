package server

import (
	"context"
	"log"
	"net"

	"github.com/Filanteyee/ADV3/proto/order"
	"github.com/Filanteyee/ADV3/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	order.UnimplementedOrderServiceServer
	orderUseCase usecase.OrderUseCase
}

func NewOrderServer(orderUseCase usecase.OrderUseCase) *OrderServer {
	return &OrderServer{
		orderUseCase: orderUseCase,
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	// Convert proto request to domain model
	orderItems := make([]usecase.OrderItem, len(req.Items))
	for i, item := range req.Items {
		orderItems[i] = usecase.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
		}
	}

	// Create order using use case
	createdOrder, err := s.orderUseCase.CreateOrder(req.UserId, orderItems)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert domain model to proto response
	protoItems := make([]*order.OrderItem, len(createdOrder.Items))
	for i, item := range createdOrder.Items {
		protoItems[i] = &order.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		}
	}

	return &order.CreateOrderResponse{
		Order: &order.Order{
			Id:          createdOrder.ID,
			UserId:      createdOrder.UserID,
			Items:       protoItems,
			TotalAmount: createdOrder.TotalAmount,
			Status:      createdOrder.Status,
			CreatedAt:   createdOrder.CreatedAt,
			UpdatedAt:   createdOrder.UpdatedAt,
		},
	}, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	order, err := s.orderUseCase.GetOrder(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoItems := make([]*order.OrderItem, len(order.Items))
	for i, item := range order.Items {
		protoItems[i] = &order.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		}
	}

	return &order.GetOrderResponse{
		Order: &order.Order{
			Id:          order.ID,
			UserId:      order.UserID,
			Items:       protoItems,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		},
	}, nil
}

func (s *OrderServer) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	orders, total, err := s.orderUseCase.ListOrders(req.UserId, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoOrders := make([]*order.Order, len(orders))
	for i, order := range orders {
		protoItems := make([]*order.OrderItem, len(order.Items))
		for j, item := range order.Items {
			protoItems[j] = &order.OrderItem{
				ProductId: item.ProductID,
				Quantity:  int32(item.Quantity),
				Price:     item.Price,
			}
		}

		protoOrders[i] = &order.Order{
			Id:          order.ID,
			UserId:      order.UserID,
			Items:       protoItems,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		}
	}

	return &order.ListOrdersResponse{
		Orders: protoOrders,
		Total:  int32(total),
	}, nil
}

func (s *OrderServer) ProcessPayment(ctx context.Context, req *order.PaymentRequest) (*order.PaymentResponse, error) {
	payment := usecase.Payment{
		OrderID:       req.OrderId,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		CardNumber:    req.CardNumber,
		CardHolder:    req.CardHolder,
		ExpiryDate:    req.ExpiryDate,
		CVV:           req.Cvv,
	}

	paymentID, status, err := s.orderUseCase.ProcessPayment(payment)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &order.PaymentResponse{
		PaymentId: paymentID,
		Status:    status,
	}, nil
}

func StartOrderServer(orderUseCase usecase.OrderUseCase, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderServer := NewOrderServer(orderUseCase)
	order.RegisterOrderServiceServer(grpcServer, orderServer)

	log.Printf("Order Service gRPC server listening on port %s", port)
	return grpcServer.Serve(lis)
}
