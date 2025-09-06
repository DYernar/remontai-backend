package service

import (
	"context"

	"github.com/DYernar/remontai-backend/internal/config"
	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/DYernar/remontai-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

type Service interface {
	LoginWithGoogle(ctx context.Context, token string, pushToken string) (domain.UserModel, error)
	LoginWithApple(ctx context.Context, token string, pushToken string) (domain.UserModel, error)
	GetUserByToken(ctx context.Context, token string) (domain.UserModel, error)
}

type service struct {
	logger *zap.SugaredLogger
	config *config.Config
	repo   postgres.Repository
}

func NewService(
	config *config.Config,
	logger *zap.SugaredLogger,
	repo postgres.Repository,
) Service {
	return &service{
		repo:   repo,
		logger: logger,
		config: config,
	}
}
