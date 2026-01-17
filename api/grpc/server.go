package grpcapi

import (
	"context"

	"github.com/chokey2nv/go-notification/domain"
	"github.com/chokey2nv/go-notification/service"
	gen "github.com/chokey2nv/go-notification/api/grpc/gen"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)
type Server struct {
	gen.UnimplementedNotificationServiceServer
	svc *service.NotificationService
}
func NewServer(svc *service.NotificationService) *Server {
	return &Server{svc: svc}
}
func (s *Server) Create(
	ctx context.Context,
	req *gen.CreateNotificationRequest,
) (*gen.Notification, error) {

	// Extract idempotency key
	var idemKey string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("idempotency-key")
		if len(values) > 0 {
			idemKey = values[0]
		}
	}

	n, err := s.svc.Add(ctx, service.AddInput{
		IdempotencyKey: idemKey,
		UserID:         req.UserId,
		Title:          req.Title,
		Message:        req.Message,
		TemplateID:     req.TemplateId,
		Channels:       mapChannels(req.Channels),
		Metadata:       req.Metadata,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return mapNotification(n), nil
}
func (s *Server) GetByID(
	ctx context.Context,
	req *gen.GetByIDRequest,
) (*gen.Notification, error) {

	n, err := s.svc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return mapNotification(n), nil
}
func (s *Server) GetByUser(
	ctx context.Context,
	req *gen.GetByUserRequest,
) (*gen.ListNotificationsResponse, error) {

	items, err := s.svc.GetByUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &gen.ListNotificationsResponse{}
	for _, n := range items {
		resp.Items = append(resp.Items, mapNotification(n))
	}

	return resp, nil
}
func mapChannels(ch []gen.Channel) []domain.Channel {
	out := make([]domain.Channel, 0, len(ch))
	for _, c := range ch {
		out = append(out, domain.Channel(c))
	}
	return out
}
func mapNotification(n *domain.Notification) *gen.Notification {
	return &gen.Notification{
		Id:        n.ID,
		UserId:    n.UserID,
		Title:     n.Title,
		Message:   n.Message,
		Channels:  mapDomainChannels(n.Channels),
		Metadata:  n.Metadata,
		CreatedAt: timestamppb.New(n.CreatedAt),
	}
}
func mapError(err error) error {
	switch err {
	case domain.ErrInvalidNotification:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
