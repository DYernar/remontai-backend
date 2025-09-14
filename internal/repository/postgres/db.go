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

	// styles
	ListStyles(ctx context.Context) ([]domain.StyleModel, error)
	GetStyleByID(ctx context.Context, styleID string) (domain.StyleModel, error)

	// image generation
	CreateImageGeneration(ctx context.Context, imageGen domain.ImageGenerationModel) (domain.ImageGenerationModel, error)
	GetImageGenerationsByUserId(ctx context.Context, userID string) ([]domain.ImageGenerationModel, error)
}

type repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) Repository {
	return &repository{
		conn: conn,
	}
}
