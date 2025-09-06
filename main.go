package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/DYernar/remontai-backend/internal/app"
	"github.com/DYernar/remontai-backend/internal/config"
	"go.uber.org/zap"
)

func init() {
	time.Local = time.UTC
}

// @title Remontai Admin API
// @version 1.0
// @description Admin API Documentation
// @host localhost:8080
// @BasePath /
func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	conf, err := config.ParseConfig()
	if err != nil {
		log.Fatal("Error while parsing config: ", err.Error())
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	sugar := logger.Sugar()

	ctx := context.Background()

	app, err := app.NewApp(ctx, conf, sugar)
	if err != nil {
		log.Fatal("Error while creating app: ", err.Error())
	}

	// return
	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
