package main

import (
	"log/slog"
	"net"
	"os"
	"user-service/databases"
	"user-service/internal/handlers"

	userV1 "github.com/ol1mov-dev/protos/pkg/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	HOST = "localhost"
	PORT = "8080"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Connecting to database...")
	DB, err := databases.ConnectPostgreSql()
	if err != nil {
		slog.Error(
			"Error connecting to database: ",
			"Details", err,
		)
	}
	slog.Info("Connected to database")

	slog.Info("Trying listening to server...", "host", HOST, "port", PORT)
	addr := net.JoinHostPort(HOST, PORT)
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		slog.Error("Failed to listen: %v", "Details", err)
	}

	slog.Info("Listening on " + addr)

	slog.Info("Starting grpc server...")
	grpcServer := grpc.NewServer()
	userV1.RegisterUserV1ServiceServer(grpcServer, &handlers.UserServerHandler{
		DB: DB,
	})

	slog.Info("GRPC server started")

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("failed to serve: ", "Details", err)
	}
}
