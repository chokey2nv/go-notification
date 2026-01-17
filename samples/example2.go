package samples

import (
	"context"
	"log"

	"github.com/chokey2nv/go-notification/dispatcher"
	"github.com/chokey2nv/go-notification/domain"
	mongorepo "github.com/chokey2nv/go-notification/repositories/mongo"
	"github.com/chokey2nv/go-notification/service"
	"github.com/chokey2nv/go-notification/worker"

	emailDispatcher "github.com/chokey2nv/go-notification/dispatcher/email"
	pushDispatcher "github.com/chokey2nv/go-notification/dispatcher/push"
	smsDispatcher "github.com/chokey2nv/go-notification/dispatcher/sms"
)

func Sample2() {
	db := connectMongo()

	notificationRepo := mongorepo.NewNotificationMongoRepository(db)
	deliveryRepo := mongorepo.NewDeliveryMongoRepository(db)
	dlqRepo := mongorepo.NewDeadLetterMongoRepository(db)
	idempotencyRepo := mongorepo.NewIdempotencyMongoRepository(db)

	smtpClient := emailDispatcher.NewDefaultSMTPClient(
		"smtp.gmail.com",
		"587",
		"no-reply@yourapp.com",
		"password",
	)

	templates, err := emailDispatcher.NewHTMLTemplateRenderer(map[string]struct {
		Subject string
		Body    string
	}{
		"welcome": {
			Subject: "Welcome {{.name}}",
			Body:    `<html><body><h2>Hello {{.name}}</h2></body></html>`,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	emailDisp := emailDispatcher.NewEmailDispatcher(
		"no-reply@yourapp.com",
		smtpClient,
		templates,
		func(ctx context.Context, userID string) (string, error) {
			return "user@email.com", nil
		},
	)

	smsClient := NewSampleSmsClient() /* Twilio / SNS implementation */

	smsDisp := smsDispatcher.NewSMSDispatcher(
		smsClient,
		templates, // reuse same templates
		func(ctx context.Context, userID string) (string, error) {
			return "+15551234567", nil
		},
	)

	pushClient := NewSamplePushClient() /* FCM / APNs implementation */

	pushDisp := pushDispatcher.NewPushDispatcher(
		pushClient,
		templates,
		func(ctx context.Context, userID string) (string, error) {
			return "device-token-123", nil
		},
	)

	txManager := mongorepo.NewMongoTransactionManager(db.Client())

	svc := service.NewNotificationService(
		notificationRepo,
		deliveryRepo,
		idempotencyRepo,
		txManager,
	)

	svc.Add(context.Background(), service.AddInput{
		IdempotencyKey: "welcome-user-123-2026-01-17", // or UUID, if failed due to network or other issues, reuse the same idempotency key
		UserID:     "user-123",
		Title:      "Welcome",
		TemplateID: "welcome",
		Channels: []domain.Channel{
			domain.ChannelEmail,
			domain.ChannelSMS,
		},
		Metadata: map[string]string{
			"name": "John",
		},
	})

	processor := worker.NewDeliveryProcessor(
		deliveryRepo,
		notificationRepo,
		dlqRepo,
		[]dispatcher.Dispatcher{
			emailDisp,
			smsDisp,
			pushDisp,
		},
	)

	processor.RunOnce(context.Background())

}
