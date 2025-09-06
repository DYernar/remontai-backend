package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Token     string `json:"token"`
	PushToken string `json:"push_token"`
	Type      string `json:"type"`
}

const (
	LoginTypeGoogle = "GOOGLE"
	LoginTypeApple  = "APPLE"
)

// LoginV1Handler godoc
// @Schemes
// @Summary Login
// @Tags auth
// @Accept			json
// @Produce		json
// @Description Login
// @Success 200 {object} domain.UserModel{}
// @Router /api/v1/auth/login [post]
// @Param book body LoginRequest true "login request"
func (a *App) LoginV1Handler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	a.logger.Info("login request", "req", req)
	if LoginTypeGoogle == req.Type {
		user, err := a.authService.LoginWithGoogle(c.Request.Context(), req.Token, req.PushToken)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"user": user})
		return
	} else if LoginTypeApple == req.Type {
		user, err := a.authService.LoginWithApple(c.Request.Context(), req.Token, req.PushToken)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"user": user})
		return
	}

	c.JSON(400, gin.H{"error": fmt.Sprintf("unknown login type: %s", req.Type)})
}
