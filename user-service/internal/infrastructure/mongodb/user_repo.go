package mongodb

import (
	"Assignment1_AbylayMoldakhmet/user-service/internal/domain"
	"Assignment1_AbylayMoldakhmet/user-service/pkg"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		collection: db.Collection("users"),
	}
}

func (r *UserRepo) Create(user *domain.User) error {
	filter := bson.M{"email": user.Email}
	if err := r.collection.FindOne(context.Background(), filter).Err(); err == nil {
		return errors.New("user already exists")
	}

	_, err := r.collection.InsertOne(context.Background(), bson.M{
		"email":    user.Email,
		"password": user.Password,
		"role":     user.Role,
	})
	return err
}

func (r *UserRepo) GetByID(id string) (*domain.User, error) {
	var user domain.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	filter := bson.M{"_id": objID}
	err = r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Update(user *domain.User) error {
	newPassword, err := pkg.HashPassword(user.Password)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"email":    user.Email,
		"password": newPassword,
	}}

	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *UserRepo) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	return err
}
