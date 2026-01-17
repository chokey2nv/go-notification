package domain

import "errors"

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidNotification  = errors.New("invalid notification")
)

func IsNotificationNotFound(err error) bool {
	return errors.Is(err, ErrNotificationNotFound)
}