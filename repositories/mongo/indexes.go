package mongorepo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func EnsureIndexes(ctx context.Context, col *mongo.Collection) error {
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.M{"user_id": 1},
		},
		{
			Keys: bson.M{"created_at": -1},
		},
	})
	return err
}
func EnsureDeliveryIndexes(ctx context.Context, col *mongo.Collection) error {
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.M{"status": 1, "next_attempt_at": 1},
		},
		{
			Keys: bson.M{"notification_id": 1},
		},
	})
	return err
}
func EnsureDLQIndexes(ctx context.Context, col *mongo.Collection) error {
	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"created_at": -1},
	})
	return err
}
