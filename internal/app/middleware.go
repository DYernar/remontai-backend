package app

import (
	"errors"

	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

var (
	testUserToken = "test"
)

type AuthedHandlerFunc func(c *gin.Context, user domain.UserModel)

func (app *App) GetUserFromHeader(c *gin.Context) (domain.UserModel, error) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		return domain.UserModel{}, errors.New("token is required")
	}

	user, err := app.authService.GetUserByToken(c, token)
	if err != nil {
		return domain.UserModel{}, err
	}

	return user, nil
}

func (app *App) PassUserMiddleware(f AuthedHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := app.GetUserFromHeader(c)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		f(c, user)
	}
}
