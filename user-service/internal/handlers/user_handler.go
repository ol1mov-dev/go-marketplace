package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	userV1 "github.com/ol1mov-dev/protos/pkg/user/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const DEFAULT_USER_ROLE = "buyer"

type UserServerHandler struct {
	userV1.UnimplementedUserV1ServiceServer
	DB *sql.DB
}

func (s *UserServerHandler) CreateUser(ctx context.Context, req *userV1.CreateUserRequest) (*userV1.CreateUserResponse, error) {
	var defaultUserRole uint32
	err := s.DB.QueryRowContext(
		ctx,
		"SELECT id FROM roles WHERE name = $1",
		DEFAULT_USER_ROLE,
	).Scan(&defaultUserRole)
	if err != nil {
		slog.Error("Failed get default role", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get user role")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.GetPassword()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		slog.Error("Failed hash password", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to process password")
	}

	var userId uint32
	err = s.DB.QueryRowContext(
		ctx,
		`
        INSERT INTO users (
            firstname,
            lastname,
            email,
            password,
            phone_number,
            role_id
        )
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
        `,
		req.GetFirstName(),
		req.GetLastName(),
		req.GetEmail(),
		string(hashedPassword),
		req.GetPhoneNumber(),
		defaultUserRole,
	).Scan(&userId)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, status.Errorf(codes.AlreadyExists, "email already exists")
		}

		slog.Error("Failed create user", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	return &userV1.CreateUserResponse{
		Id: userId,
	}, nil
}

func (s *UserServerHandler) UpdateUser(ctx context.Context, req *userV1.UpdateUserRequest) (*userV1.UpdateUserResponse, error) {

	// Хешируем пароль, если он передан
	var hashedPassword *string
	if req.GetPassword() != "" {
		hash, err := bcrypt.GenerateFromPassword(
			[]byte(req.GetPassword()),
			bcrypt.DefaultCost,
		)
		if err != nil {
			slog.Error("Failed hash password", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to process password")
		}
		hp := string(hash)
		hashedPassword = &hp
	}

	var user userV1.User

	err := s.DB.QueryRowContext(
		ctx,
		`
        UPDATE users
        SET
            email        = COALESCE(NULLIF($2, ''), email),
            password     = COALESCE($3, password),
            firstname   = COALESCE(NULLIF($4, ''), firstname),
            lastname    = COALESCE(NULLIF($5, ''), lastname),
            phone_number = COALESCE(NULLIF($6, ''), phone_number)
        WHERE id = $1
        RETURNING id, email, firstname, lastname, phone_number
        `,
		req.GetId(),
		req.GetEmail(),
		hashedPassword,
		req.GetFirstName(),
		req.GetLastName(),
		req.GetPhoneNumber(),
	).Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		slog.Error("Failed update user", "error", err, "userId", req.GetId())
		return nil, err
	}

	return &userV1.UpdateUserResponse{
		User: &user,
	}, nil
}

func (s *UserServerHandler) DeleteUser(ctx context.Context, req *userV1.DeleteUserRequest) (*userV1.DeleteUserResponse, error) {
	res, err := s.DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", req.GetId())
	if err != nil {
		slog.Error("Failed delete user", "error", err, "userId", req.GetId())
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &userV1.DeleteUserResponse{}, nil
}
