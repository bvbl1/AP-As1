package mongodb

import (
	"Assignment1_AbylayMoldakhmet/inventory-service/internal/domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepo struct {
	collection *mongo.Collection
}

func NewProductRepo(db *mongo.Database) *ProductRepo {
	return &ProductRepo{
		collection: db.Collection("products"),
	}
}

func (r *ProductRepo) Create(product *domain.Product) error {
	product.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(context.Background(), product)
	return err
}

func (r *ProductRepo) GetByID(id string) (*domain.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrProductNotFound
	}

	var product domain.Product
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepo) Update(product *domain.Product) error {
	filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": bson.M{
		"name":     product.Name,
		"price":    product.Price,
		"category": product.Category,
		"stock":    product.Stock,
	}}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *ProductRepo) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}
	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	return err
}

func (r *ProductRepo) List(filter map[string]interface{}) ([]*domain.Product, error) {
	cur, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var products []*domain.Product
	for cur.Next(context.Background()) {
		var p domain.Product
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}
