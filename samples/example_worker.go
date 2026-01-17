package samples

import (
	"context"
	"time"

	"github.com/chokey2nv/go-notification/dispatcher"
	"github.com/chokey2nv/go-notification/dispatcher/email"
	mongorepo "github.com/chokey2nv/go-notification/repositories/mongo"
	"github.com/chokey2nv/go-notification/worker"
	"go.mongodb.org/mongo-driver/mongo"
)

func WorkerExample(db *mongo.Database) {
	notificationRepo := mongorepo.NewNotificationMongoRepository(db)
	deliveryRepo := mongorepo.NewDeliveryMongoRepository(db)

	smtpClient := email.NewDefaultSMTPClient(
		"smtp.gmail.com",
		"587",
		"no-reply@yourapp.com",
		"password",
	)

	templates, _ := email.NewHTMLTemplateRenderer(map[string]struct {
		Subject string
		Body    string
	}{
		"welcome": {
			Subject: "Welcome, {{.name}}!",
			Body: `
				<!DOCTYPE html>
				<html>
				<body style="font-family: Arial, sans-serif">
					<h2>Hello {{.name}},</h2>
					<p>Welcome to <strong>{{.app}}</strong>.</p>
					<a href="{{.link}}">Get Started</a>
				</body>
				</html>
			`,
		},
	})

	emailDispatcher := email.NewEmailDispatcher(
		"no-reply@yourapp.com",
		smtpClient,
		templates,
		func(ctx context.Context, userID string) (string, error) {
			// call user service / DB
			return "user@email.com", nil
		},
	)

	// smsDispatcher := sms.NewSMSDispatcher(
	// 	smsClient,
	// 	templates,
	// 	func(ctx context.Context, userID string) (string, error) {
	// 		// call user service / DB
	// 		return "user@email.com", nil
	// 	},
	// )
	processor := worker.NewDeliveryProcessor(
		deliveryRepo,
		notificationRepo,
		nil,
		[]dispatcher.Dispatcher{
			emailDispatcher,
			// smsDispatcher,
			// pushDispatcher,
		},
	)

	for {
		_ = processor.RunOnce(context.Background())
		time.Sleep(5 * time.Second)
	}

}
