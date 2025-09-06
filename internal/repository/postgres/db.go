package postgres

import (
	"context"

	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateUser(ctx context.Context, user domain.UserModel) (domain.UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (domain.UserModel, error)
	GetUserByToken(ctx context.Context, token string) (domain.UserModel, error)
	UpsertUser(ctx context.Context, user domain.UserModel) (domain.UserModel, error)
	DeleteUser(ctx context.Context, userID string) error
	GetUserByID(ctx context.Context, userID string) (domain.UserModel, error)
}

type repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) Repository {
	return &repository{
		conn: conn,
	}
}
