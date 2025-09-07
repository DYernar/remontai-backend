package styles

import (
	"context"

	"github.com/DYernar/remontai-backend/internal/config"
	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/DYernar/remontai-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

type Service interface {
	ListStyles(ctx context.Context) ([]domain.StyleModel, error)
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

func (s *service) ListStyles(ctx context.Context) ([]domain.StyleModel, error) {
	return s.repo.ListStyles(ctx)
}
