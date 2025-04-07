package mongodb

import (
	"Assignment1_AbylayMoldakhmet/order-service/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	UpdateStatus(id string, status domain.OrderStatus) error
	GetByUserID(userID string) ([]*domain.Order, error)
}

// MongoDB реализация
type OrderRepo struct {
	collection *mongo.Collection
}

func NewOrderRepo(collection *mongo.Collection) *OrderRepo {
	return &OrderRepo{collection: collection}
}

func (r *OrderRepo) Create(order *domain.Order) error {
	_, err := r.collection.InsertOne(context.Background(), order)
	return err
}

func (r *OrderRepo) GetByID(id string) (*domain.Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var order domain.Order
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&order)
	return &order, err
}

func (r *OrderRepo) UpdateStatus(id string, status domain.OrderStatus) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

func (r *OrderRepo) GetByUserID(userID string) ([]*domain.Order, error) {
	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"user_id": userID},
	)
	if err != nil {
		return nil, err
	}

	var orders []*domain.Order
	err = cursor.All(context.Background(), &orders)
	return orders, err
}
