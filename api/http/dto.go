package httpapi

import "github.com/chokey2nv/go-notification/domain"
type CreateNotificationRequest struct {
	UserID     string            `json:"user_id"`
	Title      string            `json:"title"`
	Message    string            `json:"message"`
	TemplateID string            `json:"template_id"`
	Channels   []domain.Channel  `json:"channels"`
	Metadata   map[string]string `json:"metadata"`
}
type NotificationResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Title     string            `json:"title"`
	Message   string            `json:"message"`
	Channels  []domain.Channel  `json:"channels"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt int64             `json:"created_at"`
}
