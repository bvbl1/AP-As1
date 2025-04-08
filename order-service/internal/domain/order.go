package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID string             `bson:"user_id" json:"user_id"`
	Items  []OrderItem        `bson:"items" json:"items"`
	Status OrderStatus        `bson:"status" json:"status"`
}

type OrderItem struct {
	ProductID string  `bson:"product_id" json:"product_id"`
	Quantity  int     `bson:"quantity" json:"quantity"`
	Price     float64 `bson:"price" json:"price"`
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id string) (*Order, error)
	UpdateStatus(id string, status OrderStatus) error
	GetByUserID(userID string) ([]*Order, error)
}
