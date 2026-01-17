package mongorepo

import (
	"context"

	"github.com/chokey2nv/go-notification/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTransactionManager struct {
	client *mongo.Client
}

func NewMongoTransactionManager(
	client *mongo.Client,
) repository.TransactionManager {
	return &MongoTransactionManager{client: client}
}
func (m *MongoTransactionManager) Begin(
	ctx context.Context,
) (repository.Transaction, error) {

	session, err := m.client.StartSession()
	if err != nil {
		return nil, err
	}

	if err := session.StartTransaction(); err != nil {
		session.EndSession(ctx)
		return nil, err
	}

	txnCtx := mongo.NewSessionContext(ctx, session)

	return &MongoTransaction{
		session: session,
		ctx:     txnCtx,
	}, nil
}
