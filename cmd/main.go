package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/wanafiq/feed-api/internal/config"
	"github.com/wanafiq/feed-api/internal/database"
	"github.com/wanafiq/feed-api/internal/handlers"
	"github.com/wanafiq/feed-api/internal/logger"
	"github.com/wanafiq/feed-api/internal/middleware"
	"github.com/wanafiq/feed-api/internal/repository"
	"github.com/wanafiq/feed-api/internal/routes"
	"github.com/wanafiq/feed-api/internal/services"
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
	logger *zap.SugaredLogger

	// repositories
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	tokenRepo    repository.TokenRepository
	followerRepo repository.FollowerRepository

	// services
	authService  *services.AuthService
	emailService *services.EmailService
	userService  *services.UserService

	// handlers
	authHandler *handlers.AuthHandler
	userHandler *handlers.UserHandler

	middleware *middleware.Middleware
	router     *gin.Engine
}

func (app *application) initialize() error {
	var err error

	if app.config, err = config.LoadConfig(); err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	if app.db, err = database.InitDB(app.config.DatabaseURL); err != nil {
		return fmt.Errorf("error initializing DB: %w", err)
	}

	if app.logger, err = logger.NewLogger(); err != nil {
		return fmt.Errorf("error initializing logger: %w", err)
	}
	defer app.logger.Sync()
	fmt.Printf("logger initialized env: %s\n", app.config.Env)

	return nil
}

func main() {
	var app application
	if err := app.initialize(); err != nil {
		log.Fatalf("application initialization, %v", err)
	}
	defer app.db.Close()
	defer app.logger.Sync()

	// repositories
	app.userRepo = repository.NewUserRepository(app.db)
	app.roleRepo = repository.NewRoleRepository(app.db)
	app.tokenRepo = repository.NewTokenRepository(app.db)
	app.followerRepo = repository.NewFollowerRepository(app.db)

	// services
	app.emailService = services.NewEmailService(app.config, app.logger)
	app.authService = services.NewAuthService(
		app.config,
		app.db,
		app.logger,
		app.userRepo,
		app.roleRepo,
		app.tokenRepo,
		app.emailService,
	)
	app.userService = services.NewUserService(app.config, app.db, app.logger, app.userRepo, app.followerRepo)

	// handlers
	app.authHandler = handlers.NewAuthHandler(app.logger, app.authService)
	app.userHandler = handlers.NewUserHandler(app.logger, app.userService)

	app.middleware = middleware.NewMiddleware(app.config, app.logger)
	app.router = routes.InitRoutes(app.middleware, app.authHandler, app.userHandler)

	fmt.Printf("starting server on port %s...\n", app.config.Port)
	if err := app.router.Run(":" + app.config.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
