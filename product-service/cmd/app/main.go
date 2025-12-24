package main

import (
	"log"
	"net"
	"product-service/databases"
	"product-service/internal/handlers"

	productV1 "github.com/ol1mov-dev/protos/pkg/product/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	HOST = "localhost"
	PORT = "8081"
)

func main() {
	DB, err := databases.ConnectPostgreSql()
	if err != nil {
		log.Fatal(err)
	}

	addr := net.JoinHostPort(HOST, PORT)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	productV1.RegisterProductV1ServiceServer(grpcServer, &handlers.ProductServer{
		DB: DB,
	})

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {

		log.Fatalf("failed to serve: %v", err)
	}
}
