package app

import "github.com/gin-gonic/gin"

// ListStylesV1Handler godoc
// @Schemes
// @Summary List all styles
// @Tags styles
// @Accept json
// @Produce json
// @Description Get all available interior design styles
// @Success 200 {array} domain.StyleModel
// @Router /api/v1/styles [get]
func (a *App) ListStylesV1Handler(c *gin.Context) {
	styles, err := a.styleService.ListStyles(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"styles": styles})
}
