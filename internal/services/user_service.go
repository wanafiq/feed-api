package services

import (
	"context"
	"database/sql"
	"github.com/wanafiq/feed-api/internal/config"
	"github.com/wanafiq/feed-api/internal/models"
	"github.com/wanafiq/feed-api/internal/repository"
	"go.uber.org/zap"
)

type UserService struct {
	config       *config.Config
	db           *sql.DB
	logger       *zap.SugaredLogger
	userRepo     repository.UserRepository
	followerRepo repository.FollowerRepository
}

func NewUserService(config *config.Config, db *sql.DB, logger *zap.SugaredLogger, userRepo repository.UserRepository, followerRepo repository.FollowerRepository) *UserService {
	return &UserService{
		config:       config,
		db:           db,
		logger:       logger,
		userRepo:     userRepo,
		followerRepo: followerRepo,
	}
}

func (s *UserService) GetByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Errorw("failed to find user by id", "userID", userID, "error", err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserService) Follow(ctx context.Context, followerID string, followeeID string) error {
	_, err := s.userRepo.FindByID(ctx, followeeID)
	if err != nil {
		s.logger.Errorw("failed to find followeeID by id", "followeeID", followeeID, "error", err.Error())
		return err
	}

	if err := s.followerRepo.Save(ctx, nil, followerID, followeeID); err != nil {
		s.logger.Errorw("failed to save follower", "followerID", followerID, "followeeID", followeeID, "error", err.Error())
		return err
	}

	return nil
}

func (s *UserService) Unfollow(ctx context.Context, followerID string, followeeID string) error {
	_, err := s.userRepo.FindByID(ctx, followeeID)
	if err != nil {
		s.logger.Errorw("failed to find followeeID by id", "followeeID", followeeID, "error", err.Error())
		return err
	}

	if err := s.followerRepo.Delete(ctx, nil, followerID, followeeID); err != nil {
		s.logger.Errorw("failed to delete follower", "followerID", followerID, "followeeID", followeeID, "error", err.Error())
		return err
	}

	return nil
}

func (s *UserService) Deactivate(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Errorw("failed to find user by id", "userID", userID, "error", err.Error())
		return nil, err
	}

	user.IsActive = false

	if err := s.userRepo.Update(ctx, nil, user); err != nil {
		s.logger.Errorw("failed to update user", "userID", userID, "error", err.Error())
		return nil, err
	}

	return user, nil
}
