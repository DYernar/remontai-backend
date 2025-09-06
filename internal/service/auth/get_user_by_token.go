package auth

import (
	"context"

	"github.com/DYernar/remontai-backend/internal/domain"
)

func (s *service) GetUserByToken(ctx context.Context, token string) (domain.UserModel, error) {
	return s.repo.GetUserByToken(ctx, token)
}
