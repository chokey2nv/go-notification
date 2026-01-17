package domain

import "time"

type Notification struct {
	ID        string
	UserID    string
	Title     string
	Message   string
	Channels  []Channel
	Metadata  map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NotificationDelivery struct {
	ID             string
	NotificationID string
	Channel        Channel
	Status         DeliveryStatus
	Attempts       int
	LastError      string
	NextAttemptAt  time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
type DeadLetter struct {
	ID             string
	DeliveryID     string
	NotificationID string
	Channel        Channel
	Reason         string
	Payload        map[string]string
	CreatedAt      time.Time
}
