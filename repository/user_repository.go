package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepo{col: db.Collection("users")}
}

func (r *userRepo) Create(ctx context.Context, user interface{}) error {
	_, err := r.col.InsertOne(ctx, user)
	return err
}

func (r *userRepo) GetByID(ctx context.Context, id interface{}, result interface{}) error {
	oid, err := primitive.ObjectIDFromHex(id.(string))
	if err != nil {
		return err
	}
	return r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(result)
}

func (r *userRepo) List(ctx context.Context, filter interface{}, results interface{}) error {
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return err
	}
	return cursor.All(ctx, results)
}

func (r *userRepo) Update(ctx context.Context, id interface{}, update interface{}) error {
	oid, err := primitive.ObjectIDFromHex(id.(string))
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": update})
	return err
}

func (r *userRepo) Delete(ctx context.Context, id interface{}) error {
	oid, err := primitive.ObjectIDFromHex(id.(string))
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *userRepo) Count(ctx context.Context) (int64, error) {
	return r.col.CountDocuments(ctx, bson.D{})
}
