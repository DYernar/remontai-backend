package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/DYernar/remontai-backend/internal/util"
	goauth2 "golang.org/x/oauth2"
	googlegoauth "golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func (s *service) LoginWithGoogle(ctx context.Context, token string, pushToken string) (domain.UserModel, error) {
	userInfo, err := s.getGoogleUser(token)
	if err != nil {
		return domain.UserModel{}, fmt.Errorf("failed to get google user: %v", err)
	}

	user, err := s.repo.GetUserByEmail(ctx, userInfo.Email)

	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return domain.UserModel{}, fmt.Errorf("failed to get user by email: %v", err)
	}

	if errors.Is(err, domain.ErrUserNotFound) {
		usernameFromEmail := util.GetUsernameFromEmail(userInfo.Email)

		user = domain.UserModel{
			Email:     userInfo.Email,
			Name:      usernameFromEmail,
			Image:     userInfo.Picture,
			PushToken: pushToken,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		// create user
		user, err = s.repo.CreateUser(ctx, user)
		if err != nil {
			return domain.UserModel{}, fmt.Errorf("failed to create user: %v", err)
		}
	}

	user.Token, err = util.GenerateJWT(user.Email, user.ID)
	if err != nil {
		return domain.UserModel{}, fmt.Errorf("failed to generate jwt: %v", err)
	}

	user, err = s.repo.UpsertUser(ctx, user)
	if err != nil {
		return domain.UserModel{}, fmt.Errorf("failed to update user: %v", err)
	}

	return user, nil
}

func (s *service) getGoogleUser(accessToken string) (*oauth2.Userinfo, error) {
	ctx := context.Background()

	config := &goauth2.Config{
		ClientID:     s.config.GoogleSigninCredentials.ClientID,
		ClientSecret: s.config.GoogleSigninCredentials.ClientSecret,
		Endpoint:     googlegoauth.Endpoint,
	}

	token := &goauth2.Token{
		AccessToken: accessToken,
	}

	client := config.Client(ctx, token)

	service, err := oauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %v", err)
	}

	userInfo, err := service.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	return userInfo, nil
}
