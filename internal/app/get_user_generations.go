package app

import (
	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// GetUserGenerationsV1Handler godoc
// @Schemes
// @Summary Get user image generations
// @Tags generations
// @Accept			json
// @Produce		json
// @Description Get all image generations for current user
// @Param Authorization header string true "Paste the token here"
// @Success 200 {object} []domain.ImageGenerationModel{}
// @Router /api/v1/generations [get]
func (app *App) GetUserGenerationsV1Handler(c *gin.Context, user domain.UserModel) {
	generations, err := app.repo.GetImageGenerationsByUserId(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user generations"})
		return
	}

	c.JSON(200, generations)
}
