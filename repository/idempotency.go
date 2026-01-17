package repository

import (
	"context"

	"github.com/chokey2nv/go-notification/domain"
)

type IdempotencyRepository interface {
	Get(
		ctx context.Context,
		key string,
	) (*domain.IdempotencyRecord, error)

	Create(
		ctx context.Context,
		record *domain.IdempotencyRecord,
	) error
}
