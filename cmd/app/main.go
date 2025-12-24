package main

import (
	"log"
	"net"
	"payment-service/databases"
	"payment-service/internal/handlers"

	paymentv1 "github.com/ol1mov-dev/protos/pkg/payment/v1"
	shared1 "github.com/ol1mov-dev/shared"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	HOST = "localhost"
	PORT = "8086"
)

func main() {
	DB, err := databases.ConnectPostgreSQL()
	addr := net.JoinHostPort(HOST, PORT)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	paymentv1.RegisterPaymentServiceServer(grpcServer, &handlers.PaymentServer{
		DB: DB,
		KafkaReader: shared1.
	})

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
