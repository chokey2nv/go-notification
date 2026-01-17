package mongorepo

import (
	"context"

	"github.com/chokey2nv/go-notification/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IdempotencyMongoRepository struct {
	col *mongo.Collection
}

func NewIdempotencyMongoRepository(
	db *mongo.Database,
) *IdempotencyMongoRepository {
	return &IdempotencyMongoRepository{
		col: db.Collection("idempotency_keys"),
	}
}
func (r *IdempotencyMongoRepository) Get(
	ctx context.Context,
	key string,
) (*domain.IdempotencyRecord, error) {

	var rec domain.IdempotencyRecord
	err := r.col.FindOne(ctx, bson.M{"key": key}).Decode(&rec)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &rec, err
}
func (r *IdempotencyMongoRepository) Create(
	ctx context.Context,
	rec *domain.IdempotencyRecord,
) error {
	_, err := r.col.InsertOne(ctx, rec)
	return err
}

// index this repo
// db.idempotency_keys.createIndex({ key: 1 }, { unique: true })
