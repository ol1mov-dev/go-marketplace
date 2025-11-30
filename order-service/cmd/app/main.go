package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	orderPb "order-service/pkg/api/order_v1"
	userPb "user-servi/pkg/api/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type orderServer struct {
	orderPb.UnimplementedOrderServiceV1Server
	orders     map[string]*orderPb.Order
	mu         sync.RWMutex
	userClient userPb.UserServiceV1Client
}

func newOrderServer(userClient userPb.UserServiceV1Client) *orderServer {
	return &orderServer{
		orders:     make(map[string]*orderPb.Order),
		userClient: userClient,
	}
}

func (s *orderServer) CreateOrder(ctx context.Context, req *orderPb.CreateOrderRequest) (*orderPb.CreateOrderResponse, error) {
	log.Printf("CreateOrder called: user_id=%s, product=%s, amount=%.2f", req.UserId, req.Product, req.Amount)

	// Проверяем существование пользователя в User Service
	userResp, err := s.userClient.GetUser(ctx, &userPb.GetUserRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user not found: %v", err)
	}

	log.Printf("User verified: %s (%s)", userResp.User.Name, userResp.User.Email)

	s.mu.Lock()
	defer s.mu.Unlock()

	orderId := fmt.Sprintf("order_%d", len(s.orders)+1)
	order := &orderPb.Order{
		Id:      orderId,
		UserId:  req.UserId,
		Product: req.Product,
		Amount:  req.Amount,
	}

	s.orders[orderId] = order

	return &orderPb.CreateOrderResponse{
		Order: order,
	}, nil
}

func (s *orderServer) GetOrder(ctx context.Context, req *orderPb.GetOrderRequest) (*orderPb.GetOrderResponse, error) {
	log.Printf("GetOrder called with ID: %s", req.OrderId)

	s.mu.RLock()
	order, exists := s.orders[req.OrderId]
	s.mu.RUnlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "order not found: %s", req.OrderId)
	}

	// Получаем информацию о пользователе из User Service
	userResp, err := s.userClient.GetUser(ctx, &userPb.GetUserRequest{
		UserId: order.UserId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user info: %v", err)
	}

	return &orderPb.GetOrderResponse{
		Order:     order,
		UserName:  userResp.User.Name,
		UserEmail: userResp.User.Email,
	}, nil
}

func main() {
	// Подключаемся к User Service
	userConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	defer userConn.Close()

	userClient := userPb.NewUserServiceV1Client(userConn)

	// Запускаем Order Service
	port := ":50052"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderPb.RegisterOrderServiceV1Server(grpcServer, newOrderServer(userClient))

	log.Printf("Order Service starting on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
