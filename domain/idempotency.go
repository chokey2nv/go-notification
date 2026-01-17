package domain

import "time"

type IdempotencyRecord struct {
	Key            string
	NotificationID string
	CreatedAt      time.Time
}
