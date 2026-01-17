package dispatcher

import (
	"context"

	"github.com/chokey2nv/go-notification/domain"
)

type Dispatcher interface {
	Channel() domain.Channel
	Send(ctx context.Context, n *domain.Notification) error
}

type DispatcherFactory func(domain.Channel) Dispatcher
