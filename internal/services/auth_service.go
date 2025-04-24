package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/wanafiq/feed-api/internal/config"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/email"
	"github.com/wanafiq/feed-api/internal/models"
	"github.com/wanafiq/feed-api/internal/repository"
	"github.com/wanafiq/feed-api/internal/types"
	"github.com/wanafiq/feed-api/internal/utils"
	"go.uber.org/zap"
	"time"
)

type AuthService struct {
	config       *config.Config
	db           *sql.DB
	logger       *zap.SugaredLogger
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	tokenRepo    repository.TokenRepository
	emailService *EmailService
}

func NewAuthService(
	config *config.Config,
	db *sql.DB, logger *zap.SugaredLogger,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tokenRepo repository.TokenRepository,
	emailService *EmailService,
) *AuthService {
	return &AuthService{
		config:       config,
		db:           db,
		logger:       logger,
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		tokenRepo:    tokenRepo,
		emailService: emailService,
	}
}

func (s *AuthService) Register(ctx context.Context, req *types.RegisterRequest) (*models.User, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Errorw("userRepo.FindByEmail", "email", req.Email, "error", err.Error())
		return nil, err
	}
	if existingUser != nil {
		err := errors.New("user already exists")
		s.logger.Errorw("userRepo.FindByEmail", "email", req.Email, "error", err.Error())
		return nil, err
	}

	hashedPassword, err := utils.Hash(req.Password)
	if err != nil {
		s.logger.Errorw("utils.Hash", "error", err.Error())
		return nil, err
	}

	role, err := s.roleRepo.FindByName(ctx, constants.RoleUser)
	if err != nil {
		s.logger.Errorw("roleRepo.FindByName", "name", constants.RoleUser, "error", err.Error())
		return nil, err
	}

	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		RoleID:    role.ID,
		IsActive:  false,
		CreatedAt: time.Now(),
		CreatedBy: req.Email,
		Role:      *role,
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorw("db.BeginTx", "error", err.Error())
		return nil, err
	}
	defer tx.Rollback()

	if err := s.userRepo.Save(ctx, tx, &user); err != nil {
		s.logger.Errorw("userRepo.Save", "error", err.Error())
		return nil, err
	}

	rawToken := uuid.New().String()
	hashedToken, err := utils.Hash(rawToken)
	if err != nil {
		s.logger.Errorw("utils.Hash", "error", err.Error())
		return nil, err
	}

	token := models.Token{
		Type:      constants.ConfirmationToken,
		Value:     hashedToken,
		ExpiredAt: time.Now().Add(constants.ConfirmationTokenExpireTime),
		UserID:    user.ID,
	}

	if err := s.tokenRepo.Save(ctx, tx, &token); err != nil {
		s.logger.Errorw("tokenRepo.Save", "error", err.Error())
		return nil, err
	}

	data := types.ConfirmationEmailData{
		Username:      user.Username,
		ActivationUrl: fmt.Sprintf("%s/confirm/%s", s.config.Url.Web, rawToken),
	}

	if err := s.emailService.Send(email.ConfirmationEmail, data, &user); err != nil {
		s.logger.Errorw("emailService.Send", "error", err.Error())
		return nil, err
	}

	tx.Commit()

	return &user, nil
}

func (s *AuthService) Login(ctx context.Context, req *types.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Errorw("userRepo.FindByEmail", "error", err.Error())
		return "", err
	}

	if ok := utils.VerifyHash(user.Password, req.Password); !ok {
		s.logger.Errorw("utils.VerifyHash - invalid credential")
		return "", err
	}

	secret := s.config.Jwt.Secret
	duration := time.Duration(s.config.Jwt.ExpiryInHours) * time.Hour
	expiredAt := time.Now().Add(duration)
	issuer := s.config.Jwt.Issuer
	audience := s.config.Jwt.Audience

	token, err := utils.GenerateJWT(user, secret, expiredAt, issuer, audience)
	if err != nil {
		s.logger.Errorw("generateJWT", "error", err.Error())
		return "", err
	}

	return token, nil
}
