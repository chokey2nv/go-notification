package grpcserver

import (
	"context"

	"github.com/chokey2nv/go-notification/api/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func connectGrpc() {
	conn, _ := grpc.Dial(":9090", grpc.WithInsecure())
	client := gen.NewNotificationServiceClient(conn)

	md := metadata.New(map[string]string{
		"idempotency-key": "welcome-user-123",
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.Create(ctx, &gen.CreateNotificationRequest{
		UserId:     "user-123",
		Title:      "Welcome",
		TemplateId: "welcome",
		Channels: []gen.Channel{
			gen.Channel_CHANNEL_EMAIL,
			gen.Channel_CHANNEL_PUSH,
		},
		Metadata: map[string]string{
			"name": "John",
		},
	})

}
