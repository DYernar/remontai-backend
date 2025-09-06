package app

import (
	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// GetMeV1Handler godoc
// @Schemes
// @Summary Get current user
// @Tags auth
// @Accept			json
// @Produce		json
// @Description Get current user
// @Success 200 {object} domain.UserModel{}
// @Router /api/v1/users/me [get]
func (app *App) GetMeV1Handler(c *gin.Context, user domain.UserModel) {
	c.JSON(200, gin.H{"user": user})
}
