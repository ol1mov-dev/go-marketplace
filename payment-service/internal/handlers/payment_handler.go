package handlers

import (
	"context"
	"database/sql"
	"log"
	"time"

	paymentv1 "github.com/ol1mov-dev/protos/pkg/payment/v1"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentServer struct {
	paymentv1.UnimplementedPaymentV1ServiceServer
	DB          *sql.DB
	KafkaReader *kafka.Reader
	KafkaWriter *kafka.Writer
}

func (s *PaymentServer) CreatePayment(
	ctx context.Context,
	req *paymentv1.CreatePaymentRequest,
) (*paymentv1.CreatePaymentResponse, error) {

	query := `
        INSERT INTO payments (order_id, user_id, total_sum, status_code)
        VALUES ($1, $2, $3, $4)
        RETURNING id, status_code, created_at
    `

	var (
		id        uint32
		status    int32
		createdAt time.Time
	)

	err := s.DB.QueryRowContext(
		ctx,
		query,
		req.OrderId,
		req.UserId,
		req.TotalSum,
		paymentv1.PaymentStatus_PAYMENT_STATUS_PENDING,
	).Scan(&id, &status, &createdAt)

	if err != nil {
		return nil, err
	}

	return &paymentv1.CreatePaymentResponse{
		Id:        id,
		Status:    paymentv1.PaymentStatus(status),
		CreatedAt: timestamppb.New(createdAt),
	}, nil
}

func (s *PaymentServer) MakePayment(ctx context.Context, req *paymentv1.MakePaymentRequest) (*paymentv1.EmptyResponse, error) {
	query := "UPDATE payments SET status_code = $1 WHERE id = $2"

	_, err := s.DB.Exec(query, paymentv1.PaymentStatus_PAYMENT_STATUS_CONFIRMED, req.PaymentId)
	if err != nil {
		log.Println(err)
	}

	return &paymentv1.EmptyResponse{}, nil
}

func (s *PaymentServer) GetPayment(ctx context.Context, req *paymentv1.GetPaymentRequest) (*paymentv1.GetPaymentResponse, error) {
	return &paymentv1.GetPaymentResponse{}, nil
}
