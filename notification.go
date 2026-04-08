package main

import (
	"github.com/chokey2nv/go-notification/dispatcher"
	"github.com/chokey2nv/go-notification/repository"
	"github.com/chokey2nv/go-notification/service"
)

func New(
	repo repository.NotificationRepository,
	deliveryRepo repository.DeliveryRepository,
	idempotencyRepo repository.IdempotencyRepository,
	txManager repository.TransactionManager,
) *service.NotificationService {
	return service.NewNotificationService(repo, deliveryRepo, idempotencyRepo, txManager)
}

func NewSimple(
	repo repository.NotificationRepository,
	dispatchers []dispatcher.Dispatcher,
) *service.SimpleNotificationService {
	return service.NewSimpleNotificationService(repo, dispatchers)
}
