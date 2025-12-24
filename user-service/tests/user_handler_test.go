package tests

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"user-service/internal/handlers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	userV1 "github.com/ol1mov-dev/protos/pkg/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	return db, mock
}

func TestCreateUser_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	handler := &handlers.UserServerHandler{DB: db}

	req := &userV1.CreateUserRequest{
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john@example.com",
		Password:    "password123",
		PhoneNumber: "+1234567890",
	}

	// Mock role query
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM roles WHERE name = $1")).
		WithArgs("buyer").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta(
		"INSERT INTO users ( firstname, lastname, email, password, phone_number, role_id ) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
	)).
		WithArgs(
			req.FirstName,
			req.LastName,
			req.Email,
			sqlmock.AnyArg(), // bcrypt hash
			req.PhoneNumber,
			uint32(1), // role_id
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(100),
		)

	resp, err := handler.CreateUser(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, uint32(100), resp.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_RoleNotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	handler := &handlers.UserServerHandler{DB: db}

	req := &userV1.CreateUserRequest{
		FirstName: "John",
		Email:     "john@example.com",
		Password:  "password123",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM roles WHERE name = $1")).
		WithArgs("buyer").
		WillReturnError(sql.ErrNoRows)

	resp, err := handler.CreateUser(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	handler := &handlers.UserServerHandler{DB: db}

	req := &userV1.CreateUserRequest{
		FirstName: "John",
		Email:     "john@example.com",
		Password:  "password123",
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id FROM roles WHERE name = $1",
	)).
		WithArgs("buyer").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta(
		"INSERT INTO users ( firstname, lastname, email, password, phone_number, role_id ) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
	)).
		WithArgs(
			req.FirstName,
			"",
			req.Email,
			sqlmock.AnyArg(),
			"",
			1,
		).
		WillReturnError(&pq.Error{Code: "23505"})

	resp, err := handler.CreateUser(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Contains(t, st.Message(), "email already exists")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	handler := &handlers.UserServerHandler{DB: db}

	req := &userV1.UpdateUserRequest{
		Id:          1,
		FirstName:   "John",
		LastName:    "Updated",
		Email:       "john.updated@example.com",
		PhoneNumber: "+9876543210",
	}

	mock.ExpectQuery(`UPDATE users.*RETURNING`).
		WithArgs(
			req.Id,
			req.Email,
			nil,
			req.FirstName,
			req.LastName,
			req.PhoneNumber,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "firstname", "lastname", "phone_number"}).
			AddRow(1, req.Email, req.FirstName, req.LastName, req.PhoneNumber))

	resp, err := handler.UpdateUser(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, uint32(1), resp.User.Id)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.Equal(t, req.FirstName, resp.User.FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_WithPassword(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	handler := &handlers.UserServerHandler{DB: db}

	req := &userV1.UpdateUserRequest{
		Id:        1,
		FirstName: "John",
		Password:  "newpassword123",
	}

	mock.ExpectQuery(`UPDATE users.*RETURNING`).
		WithArgs(
			req.Id,
			"",
			sqlmock.AnyArg(), // hashed password
			req.FirstName,
			"",
			"",
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "firstname", "lastname", "phone_number"}).
			AddRow(1, "test@example.com", req.FirstName, "Doe", "+123"))

	resp, err := handler.UpdateUser(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, uint32(1), resp.User.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	handler := &handlers.UserServerHandler{DB: db}

	req := &userV1.UpdateUserRequest{
		Id:        999,
		FirstName: "John",
	}

	mock.ExpectQuery(`UPDATE users.*RETURNING`).
		WithArgs(req.Id, "", nil, req.FirstName, "", "").
		WillReturnError(sql.ErrNoRows)

	resp, err := handler.UpdateUser(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	handler := &handlers.UserServerHandler{DB: db}

	req := &userV1.UpdateUserRequest{
		Id:        1,
		FirstName: "John",
	}

	mock.ExpectQuery(`UPDATE users.*RETURNING`).
		WithArgs(req.Id, "", nil, req.FirstName, "", "").
		WillReturnError(errors.New("database error"))

	resp, err := handler.UpdateUser(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	handler := &handlers.UserServerHandler{DB: db}

	req := &userV1.DeleteUserRequest{
		Id: 1,
	}

	mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM users WHERE id = $1",
	)).
		WithArgs(req.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	resp, err := handler.DeleteUser(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_UserNotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	handler := &handlers.UserServerHandler{DB: db}

	req := &userV1.DeleteUserRequest{
		Id: 999,
	}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1")).
		WithArgs(req.Id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	resp, err := handler.DeleteUser(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	handler := &handlers.UserServerHandler{DB: db}

	req := &userV1.DeleteUserRequest{
		Id: 1,
	}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1")).
		WithArgs(req.Id).
		WillReturnError(errors.New("database error"))

	resp, err := handler.DeleteUser(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	assert.NoError(t, err)

	err = bcrypt.CompareHashAndPassword(hash, []byte("wrongpassword"))
	assert.Error(t, err)
}
