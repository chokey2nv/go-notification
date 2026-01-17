package samples

import (
	mongorepo "github.com/chokey2nv/go-notification/repositories/mongo"
	"github.com/chokey2nv/go-notification/service"
)

func DefaultNotificationServer() *service.NotificationService {
	db := connectMongo()

	notificationRepo := mongorepo.NewNotificationMongoRepository(db)
	deliveryRepo := mongorepo.NewDeliveryMongoRepository(db)
	idempotencyRepo := mongorepo.NewIdempotencyMongoRepository(db)

	txManager := mongorepo.NewMongoTransactionManager(db.Client())

	svc := service.NewNotificationService(
		notificationRepo,
		deliveryRepo,
		idempotencyRepo,
		txManager,
	)

	return svc

}
