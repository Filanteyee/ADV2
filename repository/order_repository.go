package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Filanteyee/ADV3/usecase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository interface {
	SaveOrder(order *usecase.Order) error
	GetOrderByID(orderID string) (*usecase.Order, error)
	GetOrdersByUserID(userID string, page, limit int) ([]usecase.Order, int, error)
	UpdateOrder(order *usecase.Order) error
}

type orderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository() OrderRepository {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	db := client.Database("ecommerce")
	return &orderRepository{
		collection: db.Collection("orders"),
	}
}

func (r *orderRepository) SaveOrder(order *usecase.Order) error {
	_, err := r.collection.InsertOne(context.Background(), order)
	return err
}

func (r *orderRepository) GetOrderByID(orderID string) (*usecase.Order, error) {
	var order usecase.Order
	err := r.collection.FindOne(context.Background(), bson.M{"id": orderID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetOrdersByUserID(userID string, page, limit int) ([]usecase.Order, int, error) {
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{"createdat", -1}})

	cursor, err := r.collection.Find(context.Background(), bson.M{"userid": userID}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var orders []usecase.Order
	if err = cursor.All(context.Background(), &orders); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(context.Background(), bson.M{"userid": userID})
	if err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

func (r *orderRepository) UpdateOrder(order *usecase.Order) error {
	order.UpdatedAt = time.Now().Format(time.RFC3339)
	_, err := r.collection.ReplaceOne(context.Background(), bson.M{"id": order.ID}, order)
	return err
}
