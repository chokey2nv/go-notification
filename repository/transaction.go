package repository

import "context"

type Transaction interface {
	Context() context.Context
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TransactionManager interface {
	Begin(ctx context.Context) (Transaction, error)
}
