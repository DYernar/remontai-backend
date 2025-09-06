package app

import "github.com/gin-gonic/gin"

// @Remontai Ping endpoint
// @Description Returns pong message
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/ping [get]
func (a *App) PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
