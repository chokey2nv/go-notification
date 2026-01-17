package repository

import (
	"context"

	"github.com/chokey2nv/go-notification/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *domain.Notification) error
	Update(ctx context.Context, n *domain.Notification) error
	GetByID(ctx context.Context, id string) (*domain.Notification, error)
	GetByUser(ctx context.Context, userID string) ([]*domain.Notification, error)
}

type DeliveryRepository interface {
	Create(ctx context.Context, d *domain.NotificationDelivery) error
	Update(ctx context.Context, d *domain.NotificationDelivery) error
	GetPending(ctx context.Context, limit int) ([]*domain.NotificationDelivery, error)
	GetByID(ctx context.Context, id string) (*domain.NotificationDelivery, error)
	Claim(
		ctx context.Context,
		id string,
	) (*domain.NotificationDelivery, error)
}

type DeadLetterRepository interface {
	Create(ctx context.Context, d *domain.DeadLetter) error
	Get(ctx context.Context, limit int) ([]*domain.DeadLetter, error)
	GetByID(ctx context.Context, id string) (*domain.DeadLetter, error)
	Delete(ctx context.Context, id string) error
}
