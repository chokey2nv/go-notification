package grpcserver

import (
	"log"
	"net"

	grpcapi "github.com/chokey2nv/go-notification/api/grpc"
	"github.com/chokey2nv/go-notification/api/grpc/gen"
	"github.com/chokey2nv/go-notification/samples"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	svc := samples.DefaultNotificationServer()

	gen.RegisterNotificationServiceServer(
		grpcServer,
		grpcapi.NewServer(svc),
	)

	log.Println("gRPC server listening on :9090")
	grpcServer.Serve(lis)

}
