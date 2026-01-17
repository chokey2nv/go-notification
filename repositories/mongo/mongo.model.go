package mongorepo

import (
	"time"

	"github.com/chokey2nv/go-notification/domain"
)

type notificationDocument struct {
	ID        string            `bson:"_id"`
	UserID    string            `bson:"user_id"`
	Title     string            `bson:"title"`
	Message   string            `bson:"message"`
	Channels  []domain.Channel  `bson:"channels"`
	Metadata  map[string]string `bson:"metadata"`
	CreatedAt time.Time         `bson:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at"`
}

type deliveryDocument struct {
	ID             string                `bson:"_id"`
	NotificationID string                `bson:"notification_id"`
	Channel        domain.Channel        `bson:"channel"`
	Status         domain.DeliveryStatus `bson:"status"`
	Attempts       int                   `bson:"attempts"`
	LastError      string                `bson:"last_error,omitempty"`
	NextAttemptAt  time.Time             `bson:"next_attempt_at"`
	CreatedAt      time.Time             `bson:"created_at"`
	UpdatedAt      time.Time             `bson:"updated_at"`
}

type deadLetterDocument struct {
	ID             string            `bson:"_id"`
	DeliveryID     string            `bson:"delivery_id"`
	NotificationID string            `bson:"notification_id"`
	Channel        domain.Channel    `bson:"channel"`
	Reason         string            `bson:"reason"`
	Payload        map[string]string `bson:"payload"`
	CreatedAt      time.Time         `bson:"created_at"`
}
