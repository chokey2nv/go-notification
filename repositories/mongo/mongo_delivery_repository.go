package mongorepo

import (
	"context"
	"time"

	"github.com/chokey2nv/go-notification/domain"
	"github.com/chokey2nv/go-notification/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeliveryMongoRepository struct {
	collection *mongo.Collection
}

func NewDeliveryMongoRepository(
	db *mongo.Database,
) repository.DeliveryRepository {
	return &DeliveryMongoRepository{
		collection: db.Collection("notification_deliveries"),
	}
}
func (r *DeliveryMongoRepository) Create(
	ctx context.Context,
	d *domain.NotificationDelivery,
) error {

	doc := deliveryDocument{
		ID:             d.ID,
		NotificationID: d.NotificationID,
		Channel:        d.Channel,
		Status:         d.Status,
		Attempts:       d.Attempts,
		LastError:      d.LastError,
		NextAttemptAt:  d.NextAttemptAt,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	return err
}
func (r *DeliveryMongoRepository) Update(
	ctx context.Context,
	d *domain.NotificationDelivery,
) error {

	filter := bson.M{
		"_id": d.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"status":          d.Status,
			"attempts":        d.Attempts,
			"last_error":      d.LastError,
			"next_attempt_at": d.NextAttemptAt,
			"updated_at":      time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
func (r *DeliveryMongoRepository) GetPending(
	ctx context.Context,
	limit int,
) ([]*domain.NotificationDelivery, error) {

	filter := bson.M{
		"status": bson.M{
			"$in": []domain.DeliveryStatus{
				domain.DeliveryPending,
				domain.DeliveryRetrying,
			},
		},
		"next_attempt_at": bson.M{"$lte": time.Now()},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"next_attempt_at": 1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*domain.NotificationDelivery

	for cursor.Next(ctx) {
		var doc deliveryDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		result = append(result, mapDeliveryToDomain(&doc))
	}

	return result, nil
}
func (r *DeliveryMongoRepository) GetByID(
	ctx context.Context,
	id string,
) (*domain.NotificationDelivery, error) {

	var doc deliveryDocument

	err := r.collection.
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&doc)

	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotificationNotFound
	}
	if err != nil {
		return nil, err
	}

	return mapDeliveryToDomain(&doc), nil
}


////////////////////////
// Worker
////////////////////////
func (r *DeliveryMongoRepository) Claim(
	ctx context.Context,
	id string,
) (*domain.NotificationDelivery, error) {

	filter := bson.M{
		"_id":    id,
		"status": bson.M{"$in": []domain.DeliveryStatus{
			domain.DeliveryPending,
			domain.DeliveryRetrying,
		}},
	}

	update := bson.M{
		"$set": bson.M{
			"status":     domain.DeliveryRetrying,
			"updated_at": time.Now(),
		},
	}

	var doc deliveryDocument
	err := r.collection.
		FindOneAndUpdate(ctx, filter, update).
		Decode(&doc)

	if err == mongo.ErrNoDocuments {
		return nil, nil // already claimed
	}
	if err != nil {
		return nil, err
	}

	return mapDeliveryToDomain(&doc), nil
}



func mapDeliveryToDomain(doc *deliveryDocument) *domain.NotificationDelivery {
	return &domain.NotificationDelivery{
		ID:             doc.ID,
		NotificationID: doc.NotificationID,
		Channel:        doc.Channel,
		Status:         doc.Status,
		Attempts:       doc.Attempts,
		LastError:      doc.LastError,
		NextAttemptAt:  doc.NextAttemptAt,
		CreatedAt:      doc.CreatedAt,
		UpdatedAt:      doc.UpdatedAt,
	}
}
