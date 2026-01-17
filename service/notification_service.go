package service

import (
	"context"
	"time"

	"github.com/chokey2nv/go-notification/domain"
	"github.com/chokey2nv/go-notification/repository"
	"github.com/google/uuid"
)

type NotificationService struct {
	notifications repository.NotificationRepository
	deliveries    repository.DeliveryRepository
	idempotency   repository.IdempotencyRepository
	txManager     repository.TransactionManager
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	deliveryRepo repository.DeliveryRepository,
	idempotencyRepo repository.IdempotencyRepository,
	txManager repository.TransactionManager,
) *NotificationService {
	return &NotificationService{
		notifications: notificationRepo,
		deliveries:    deliveryRepo,
		idempotency:   idempotencyRepo,
		txManager:     txManager,
	}
}

type AddInput struct {
	IdempotencyKey string // optional but strongly recommended
	UserID         string
	Title          string
	Message        string
	TemplateID     string
	Channels       []domain.Channel
	Metadata       map[string]string
}

/* func (s *NotificationService) Add(
	ctx context.Context,
	in AddInput,
) (*domain.Notification, error) {

	if in.UserID == "" {
		return nil, domain.ErrInvalidNotification
	}
	if len(in.Channels) == 0 {
		return nil, domain.ErrInvalidNotification
	}
	if in.Message == "" && in.TemplateID == "" {
		return nil, domain.ErrInvalidNotification
	}

	meta := in.Metadata
	if meta == nil {
		meta = make(map[string]string)
	}
	if in.TemplateID != "" {
		meta["template_id"] = in.TemplateID
	}

	now := time.Now()

	n := &domain.Notification{
		ID:        uuid.NewString(),
		UserID:    in.UserID,
		Title:     in.Title,
		Message:   in.Message,
		Channels:  in.Channels,
		Metadata:  meta,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 1️⃣ Persist notification
	if err := s.notifications.Create(ctx, n); err != nil {
		return nil, err
	}

	// 2️⃣ Create delivery records (ASYNC PATH)
	for _, ch := range in.Channels {
		delivery := &domain.NotificationDelivery{
			ID:             uuid.NewString(),
			NotificationID: n.ID,
			Channel:        ch,
			Status:         domain.DeliveryPending,
			Attempts:       0,
			NextAttemptAt:  now,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if err := s.deliveries.Create(ctx, delivery); err != nil {
			return nil, err
		}
	}

	return n, nil
} */

func (s *NotificationService) Add(
	ctx context.Context,
	in AddInput,
) (*domain.Notification, error) {

	if in.UserID == "" || len(in.Channels) == 0 {
		return nil, domain.ErrInvalidNotification
	}
	if in.Message == "" && in.TemplateID == "" {
		return nil, domain.ErrInvalidNotification
	}

	// 1️⃣ Fast-path: idempotency check
	if in.IdempotencyKey != "" {
		existing, err := s.idempotency.Get(ctx, in.IdempotencyKey)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return s.notifications.GetByID(
				ctx,
				existing.NotificationID,
			)
		}
	}

	meta := in.Metadata
	if meta == nil {
		meta = make(map[string]string)
	}
	if in.TemplateID != "" {
		meta["template_id"] = in.TemplateID
	}

	now := time.Now()

	// 2️⃣ Begin transaction
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	tctx := tx.Context()

	n := &domain.Notification{
		ID:        uuid.NewString(),
		UserID:    in.UserID,
		Title:     in.Title,
		Message:   in.Message,
		Channels:  in.Channels,
		Metadata:  meta,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 3️⃣ Create notification
	if err := s.notifications.Create(tctx, n); err != nil {
		return nil, err
	}

	// 4️⃣ Create deliveries
	for _, ch := range in.Channels {
		d := &domain.NotificationDelivery{
			ID:             uuid.NewString(),
			NotificationID: n.ID,
			Channel:        ch,
			Status:         domain.DeliveryPending,
			NextAttemptAt:  now,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if err := s.deliveries.Create(tctx, d); err != nil {
			return nil, err
		}
	}

	// 5️⃣ Record idempotency key (inside same txn)
	if in.IdempotencyKey != "" {
		rec := &domain.IdempotencyRecord{
			Key:            in.IdempotencyKey,
			NotificationID: n.ID,
			CreatedAt:      now,
		}
		if err := s.idempotency.Create(tctx, rec); err != nil {
			return nil, err
		}
	}

	// 6️⃣ Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return n, nil
}


func (s *NotificationService) Update(
	ctx context.Context,
	n *domain.Notification,
) error {
	n.UpdatedAt = time.Now()
	return s.notifications.Update(ctx, n)
}
func (s *NotificationService) GetByID(
	ctx context.Context,
	id string,
) (*domain.Notification, error) {
	return s.notifications.GetByID(ctx, id)
}

func (s *NotificationService) GetByUser(
	ctx context.Context,
	userID string,
) ([]*domain.Notification, error) {
	return s.notifications.GetByUser(ctx, userID)
}
