package samples

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func Sample(db *mongo.Database) {
	// repo := mongorepo.NewNotificationMongoRepository(db)

	// smtpClient := email.NewDefaultSMTPClient(
	// 	"smtp.gmail.com",
	// 	"587",
	// 	"no-reply@yourapp.com",
	// 	"password",
	// )

	// templates, _ := email.NewHTMLTemplateRenderer(map[string]struct {
	// 	Subject string
	// 	Body    string
	// }{
	// 	"welcome": {
	// 		Subject: "Welcome, {{.name}}!",
	// 		Body: `
	// 			<!DOCTYPE html>
	// 			<html>
	// 			<body style="font-family: Arial, sans-serif">
	// 				<h2>Hello {{.name}},</h2>
	// 				<p>Welcome to <strong>{{.app}}</strong>.</p>
	// 				<a href="{{.link}}">Get Started</a>
	// 			</body>
	// 			</html>
	// 		`,
	// 	},
	// })

	// emailDispatcher := email.NewEmailDispatcher(
	// 	"no-reply@yourapp.com",
	// 	smtpClient,
	// 	templates,
	// 	func(ctx context.Context, userID string) (string, error) {
	// 		// call user service / DB
	// 		return "user@email.com", nil
	// 	},
	// )
	// sms.SMSDispatcher
	// sms.NewSMSDispatcher()

	// svc := notification.NewSimple(
	// 	repo,
	// 	[]dispatcher.Dispatcher{
	// 		emailDispatcher,
	// 	},
	// )
	// svc.Add(context.Background(), service.SimpleAddInput{
	// 	UserID: "123",
	// 	Title:  "Welcome",
	// 	// Message: "Welcome, {{.name}}!",
	// 	TemplateID: "welcome",
	// 	Channels:   []domain.Channel{domain.ChannelEmail},
	// 	Metadata: map[string]string{
	// 		"name": "John",
	// 		"app":  "MyApp",
	// 		"link": "https://myapp.com/start",
	// 	},
	// })
	// create sample svc.Add
	// svc.Add(
	// 	context.Background(),
	// 	"123",
	// 	"Welcome",
	// 	"", // "Welcome, {{.name}}!",
	// 	[]domain.Channel{domain.ChannelEmail},
	// 	map[string]string{"name": "John"},
	// )
}
