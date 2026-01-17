package service

import (
	"context"
	"time"

	"github.com/chokey2nv/go-notification/dispatcher"
	"github.com/chokey2nv/go-notification/domain"
	"github.com/chokey2nv/go-notification/repository"
	"github.com/google/uuid"
)

type SimpleNotificationService struct {
	repo        repository.NotificationRepository
	dispatchers map[domain.Channel]dispatcher.Dispatcher
}

func NewSimpleNotificationService(
	repo repository.NotificationRepository,
	dispatchers []dispatcher.Dispatcher,
) *SimpleNotificationService {
	dMap := make(map[domain.Channel]dispatcher.Dispatcher)
	for _, d := range dispatchers {
		dMap[d.Channel()] = d
	}

	return &SimpleNotificationService{
		repo:        repo,
		dispatchers: dMap,
	}
}

type SimpleAddInput struct {
	UserID     string
	Title      string
	Message    string // optional
	TemplateID string // optional
	Channels   []domain.Channel
	Metadata   map[string]string
}

func (s *SimpleNotificationService) Add(
	ctx context.Context,
	in SimpleAddInput,
) (*domain.Notification, error) {
	userID := in.UserID
	title := in.Title
	message := in.Message
	channels := in.Channels
	meta := in.Metadata

	if in.TemplateID != "" {
		meta["template_id"] = in.TemplateID
	}

	n := &domain.Notification{
		ID:        uuid.NewString(),
		UserID:    userID,
		Title:     title,
		Message:   message,
		Channels:  channels,
		Metadata:  meta,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}

	// Fire-and-forget or sync (your choice later)
	for _, ch := range channels {
		if d, ok := s.dispatchers[ch]; ok {
			_ = d.Send(ctx, n) // intentionally not blocking core flow
		}
	}

	return n, nil
}

func (s *SimpleNotificationService) Update(
	ctx context.Context,
	n *domain.Notification,
) error {
	n.UpdatedAt = time.Now()
	return s.repo.Update(ctx, n)
}
func (s *SimpleNotificationService) GetByID(
	ctx context.Context,
	id string,
) (*domain.Notification, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SimpleNotificationService) GetByUser(
	ctx context.Context,
	userID string,
) ([]*domain.Notification, error) {
	return s.repo.GetByUser(ctx, userID)
}
