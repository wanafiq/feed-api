package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/wanafiq/feed-api/internal/config"
	"github.com/wanafiq/feed-api/internal/database"
	"github.com/wanafiq/feed-api/internal/logger"
	"github.com/wanafiq/feed-api/internal/routes"
	"go.uber.org/zap"
	"log"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Failed to load .env:", err)
	}
}

type application struct {
	config *config.Config
	db     *sql.DB
	logger *zap.Logger
	router *gin.Engine
}

func (app *application) initialize() error {
	var err error

	if app.config, err = config.LoadConfig(); err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	if app.db, err = database.InitDB(); err != nil {
		return fmt.Errorf("error initializing DB: %w", err)
	}

	if app.logger, err = logger.NewLogger(); err != nil {
		return fmt.Errorf("error initializing logger: %w", err)
	}
	fmt.Printf("Logger initialized env: %s\n", app.config.Env)

	app.router = routes.InitRoutes()

	return nil
}

func main() {
	var app application
	if err := app.initialize(); err != nil {
		log.Fatalf("Application initialization, %v", err)
	}
	defer app.db.Close()
	defer app.logger.Sync()

	fmt.Printf("Starting server on port %s...\n", app.config.Port)
	if err := app.router.Run(":" + app.config.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
