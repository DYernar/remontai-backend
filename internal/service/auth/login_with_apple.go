package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/DYernar/remontai-backend/internal/util"
	"github.com/Timothylock/go-signin-with-apple/apple"
)

func (s *service) LoginWithApple(ctx context.Context, token string, pushToken string) (domain.UserModel, error) {
	email, err := s.getAppleUser(token)
	if err != nil {
		s.logger.Errorf("failed to get apple user %s", err)
		return domain.UserModel{}, fmt.Errorf("failed to get apple user: %v", err)
	}

	user, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		s.logger.Error("failed to get user by email", "error", err)
		return domain.UserModel{}, fmt.Errorf("failed to get user by email: %v", err)
	}

	if errors.Is(err, domain.ErrUserNotFound) {
		usernameFromEmail := util.GetUsernameFromEmail(email)

		user = domain.UserModel{
			Email:     email,
			Name:      usernameFromEmail,
			Image:     "",
			PushToken: pushToken,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		// create user
		user, err = s.repo.CreateUser(ctx, user)
		if err != nil {
			s.logger.Error("failed to create user", "error", err)
			return domain.UserModel{}, fmt.Errorf("failed to create user: %v", err)
		}
	}

	user.Token, err = util.GenerateJWT(user.Email, user.ID)
	if err != nil {
		s.logger.Error("failed to generate jwt", "error", err)
		return domain.UserModel{}, fmt.Errorf("failed to generate jwt: %v", err)
	}

	user, err = s.repo.UpsertUser(ctx, user)
	if err != nil {
		s.logger.Error("failed to update user", "error", err)
		return domain.UserModel{}, fmt.Errorf("failed to update user: %v", err)
	}

	return user, nil
}

func (s *service) getAppleUser(token string) (string, error) {
	secret, err := apple.GenerateClientSecret(
		s.config.AppleSigninCredentials.PrivateKey,
		s.config.AppleSigninCredentials.TeamID,
		s.config.AppleSigninCredentials.ClientID,
		s.config.AppleSigninCredentials.KeyID,
	)
	if err != nil {
		s.logger.Error("failed to generate apple client secret %s", err)
		return "", err
	}
	client := apple.New()

	vReq := apple.AppValidationTokenRequest{
		ClientID:     s.config.AppleSigninCredentials.ClientID,
		ClientSecret: secret,
		Code:         token,
	}

	var resp apple.ValidationResponse

	// Do the verification
	err = client.VerifyAppToken(context.Background(), vReq, &resp)
	if err != nil {
		s.logger.Errorf("failed to verify apple token %s", err)
		return "", err
	}
	if resp.Error != "" {
		s.logger.Errorf("failed to verify apple token %s", resp.Error)
		return "", errors.New(resp.Error)
	}

	// Get the email
	claim, err := apple.GetClaims(resp.IDToken)
	if err != nil {
		s.logger.Error("failed to get apple claims %s", err)
		return "", err
	}
	email := (*claim)["email"]

	return fmt.Sprintf("%s", email), nil
}
