package mongorepo

import (
	"context"

	"github.com/chokey2nv/go-notification/domain"
	"github.com/chokey2nv/go-notification/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeadLetterMongoRepository struct {
	collection *mongo.Collection
}

// Delete implements repository.DeadLetterRepository.
func (r *DeadLetterMongoRepository) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetByID implements repository.DeadLetterRepository.
func (r *DeadLetterMongoRepository) GetByID(ctx context.Context, id string) (*domain.DeadLetter, error) {
	panic("unimplemented")
}

func NewDeadLetterMongoRepository(
	db *mongo.Database,
) repository.DeadLetterRepository {
	return &DeadLetterMongoRepository{
		collection: db.Collection("notification_dlq"),
	}
}
func (r *DeadLetterMongoRepository) Create(
	ctx context.Context,
	d *domain.DeadLetter,
) error {

	doc := deadLetterDocument{
		ID:             d.ID,
		DeliveryID:     d.DeliveryID,
		NotificationID: d.NotificationID,
		Channel:        d.Channel,
		Reason:         d.Reason,
		Payload:        d.Payload,
		CreatedAt:      d.CreatedAt,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	return err
}
func (r *DeadLetterMongoRepository) Get(
	ctx context.Context,
	limit int,
) ([]*domain.DeadLetter, error) {

	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*domain.DeadLetter

	for cursor.Next(ctx) {
		var doc deadLetterDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		result = append(result, mapDLQToDomain(&doc))
	}

	return result, nil
}
func mapDLQToDomain(doc *deadLetterDocument) *domain.DeadLetter {
	return &domain.DeadLetter{
		ID:             doc.ID,
		DeliveryID:     doc.DeliveryID,
		NotificationID: doc.NotificationID,
		Channel:        doc.Channel,
		Reason:         doc.Reason,
		Payload:        doc.Payload,
		CreatedAt:      doc.CreatedAt,
	}
}
