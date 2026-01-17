package mongorepo

import (
	"context"

	"github.com/chokey2nv/go-notification/domain"
	"github.com/chokey2nv/go-notification/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationMongoRepository struct {
	collection *mongo.Collection
}

func NewNotificationMongoRepository(
	db *mongo.Database,
) repository.NotificationRepository {
	return &NotificationMongoRepository{
		collection: db.Collection("notifications"),
	}
}
func (r *NotificationMongoRepository) Create(
	ctx context.Context,
	n *domain.Notification,
) error {
	doc := notificationDocument{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Message:   n.Message,
		Channels:  n.Channels,
		Metadata:  n.Metadata,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	return err
}
func (r *NotificationMongoRepository) Update(
	ctx context.Context,
	n *domain.Notification,
) error {
	filter := bson.M{"_id": n.ID}
	update := bson.M{
		"$set": bson.M{
			"title":      n.Title,
			"message":    n.Message,
			"channels":   n.Channels,
			"metadata":   n.Metadata,
			"updated_at": n.UpdatedAt,
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return domain.ErrNotificationNotFound
	}

	return nil
}
func (r *NotificationMongoRepository) GetByID(
	ctx context.Context,
	id string,
) (*domain.Notification, error) {

	var doc notificationDocument

	err := r.collection.
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&doc)

	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotificationNotFound
	}
	if err != nil {
		return nil, err
	}

	return mapToDomain(&doc), nil
}
func (r *NotificationMongoRepository) GetByUser(
	ctx context.Context,
	userID string,
) ([]*domain.Notification, error) {

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*domain.Notification

	for cursor.Next(ctx) {
		var doc notificationDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		result = append(result, mapToDomain(&doc))
	}

	return result, nil
}
func mapToDomain(doc *notificationDocument) *domain.Notification {
	return &domain.Notification{
		ID:        doc.ID,
		UserID:    doc.UserID,
		Title:     doc.Title,
		Message:   doc.Message,
		Channels:  doc.Channels,
		Metadata:  doc.Metadata,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}
