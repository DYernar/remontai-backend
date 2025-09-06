package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/DYernar/remontai-backend/docs"
	"github.com/DYernar/remontai-backend/internal/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

const (
	DB_ROOT_PASSWORD = "DB_ROOT_PASSWORD"
	ENV_HOST         = "HOST"
)

type App struct {
	config      *config.Config
	logger      *zap.SugaredLogger
	Router      *gin.Engine
	SwaggerHost string
	Host        string
	Port        string
	DBDSN       string
}

func NewApp(
	ctx context.Context,
	config *config.Config,
	logger *zap.SugaredLogger,
) (*App, error) {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	swaggerHost := os.Getenv("SERVICE_URL")
	appHost := swaggerHost
	if swaggerHost == "" {
		swaggerHost = "localhost:" + port
		appHost = config.AppConfig.Host
	}

	logger.Infof("Swagger host: %s", swaggerHost)

	swaggerHost = strings.ReplaceAll(swaggerHost, "https://", "")
	swaggerHost = strings.ReplaceAll(swaggerHost, "http://", "")

	app := &App{
		config:      config,
		logger:      logger,
		SwaggerHost: swaggerHost,
		Host:        appHost,
		Port:        port,
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	gin.SetMode(gin.ReleaseMode)

	schemes := []string{"http"}

	if !strings.Contains(a.SwaggerHost, "localhost") {
		schemes = []string{"https"}
	}

	a.Router = gin.New()

	a.Router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	docs.SwaggerInfo.Title = "Remontai Admin API"
	docs.SwaggerInfo.Description = "Admin API Documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = a.SwaggerHost
	docs.SwaggerInfo.Schemes = schemes
	docs.SwaggerInfo.BasePath = "/"

	a.Router.Use(gin.Recovery())

	// use default gin logger
	a.Router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - %s - %s - %s - %d %s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.ClientIP,
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	}))

	a.Router.MaxMultipartMemory = 8 << 20 // 8 MiB
	// swagger docs
	a.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	a.Router.GET("/api/v1/ping", a.PingHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	serverErrors := make(chan error, 1)

	go func() {
		a.logger.Infof("Starting server on %s", a.Host)
		if err := a.Router.Run(":" + a.Port); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		a.logger.Fatalf("Error starting server: %v", err)
		return err
	case <-quit:
		a.logger.Info("Shutting down server...")
	}

	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return nil
}
