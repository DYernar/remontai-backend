package imagegeneration

import (
	"context"
	"mime/multipart"

	"github.com/DYernar/remontai-backend/internal/config"
	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/DYernar/remontai-backend/internal/repository/postgres"
	"github.com/DYernar/remontai-backend/internal/repository/s3"
	"go.uber.org/zap"
)

type Service interface {
	GetImageGenerationsByUserId(ctx context.Context, userID string) ([]domain.ImageGenerationModel, error)
	QuickGenerateImage(
		ctx context.Context,
		userID string,
		imageFile multipart.File,
		imageHeader *multipart.FileHeader,
		roomType,
		styleID string,
	) (domain.ImageGenerationModel, error)
}

type service struct {
	logger *zap.SugaredLogger
	config *config.Config
	s3     s3.S3Repository
	repo   postgres.Repository
}

func NewService(
	config *config.Config,
	logger *zap.SugaredLogger,
	s3 s3.S3Repository,
	repo postgres.Repository,
) Service {
	return &service{
		repo:   repo,
		logger: logger,
		s3:     s3,
		config: config,
	}
}

func (s *service) GetImageGenerationsByUserId(ctx context.Context, userID string) ([]domain.ImageGenerationModel, error) {
	return s.repo.GetImageGenerationsByUserId(ctx, userID)
}
