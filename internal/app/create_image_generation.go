package app

import (
	"github.com/DYernar/remontai-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// CreateImageGenerationV1Handler godoc
// @Schemes
// @Summary Create image generation
// @Tags generations
// @Accept			multipart/form-data
// @Produce		json
// @Description Create a new image generation
// @Param roomtype formData string true "Room type"
// @Param styleid formData string false "Style ID"
// @Param image formData file true "Image file"
// @Success 200 {object} domain.ImageGenerationModel{}
// @Param Authorization header string true "Paste the token here"
// @Router /api/v1/generations [post]
func (app *App) CreateImageGenerationV1Handler(c *gin.Context, user domain.UserModel) {
	// Parse multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Get form values
	roomType := c.PostForm("roomtype")
	if roomType == "" {
		c.JSON(400, gin.H{"error": "roomtype is required"})
		return
	}

	styleID := c.PostForm("styleid")

	// Get uploaded file
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	resp, err := app.generateService.QuickGenerateImage(c.Request.Context(), user.ID, file, header, roomType, styleID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"generation": resp})
}
