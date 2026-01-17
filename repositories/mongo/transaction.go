package mongorepo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTransaction struct {
	session mongo.Session
	ctx     context.Context
}

func (t *MongoTransaction) Context() context.Context {
	return t.ctx
}

func (t *MongoTransaction) Commit(ctx context.Context) error {
	return t.session.CommitTransaction(ctx)
}

func (t *MongoTransaction) Rollback(ctx context.Context) error {
	return t.session.AbortTransaction(ctx)
}
