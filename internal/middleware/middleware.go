package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wanafiq/feed-api/internal/config"
	"github.com/wanafiq/feed-api/internal/utils"
	"go.uber.org/zap"
	"net/http"
	"slices"
)

const (
	authHeaderKey  = "Authorization"
	UserContextKey = "userCtx"
)

type Middleware struct {
	config *config.Config
	logger *zap.SugaredLogger
}

func NewMiddleware(config *config.Config, logger *zap.SugaredLogger) *Middleware {
	return &Middleware{
		config: config,
		logger: logger,
	}
}

type UserContext struct {
	ID       string
	Username string
	Email    string
	IsActive bool
	Role     string
}

func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authHeaderKey)

		claims, err := utils.ParseAndValidateJWT(authHeader, m.config.Jwt.Secret)
		if err != nil {
			m.logger.Errorw("failed to parse JWT", "error", err, "authHeader", authHeader)
			m.abortWithJSON(c, http.StatusUnauthorized, err.Error())
			return
		}

		userCtx := UserContext{
			ID:   claims.Subject,
			Role: claims.Role,
		}

		c.Set(UserContextKey, userCtx)

		c.Next()
	}
}

func (m *Middleware) RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx, exists := GetUserContext(c)
		if !exists {
			m.logger.Errorw("GetUserContext", "exists", false)
			m.abortWithJSON(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		userRole := userCtx.Role
		authorized := slices.Contains(allowedRoles, userRole)
		if !authorized {
			m.abortWithJSON(c, http.StatusForbidden, "insufficient permissions")
			return
		}

		c.Next()
	}
}

func GetUserContext(c *gin.Context) (UserContext, bool) {
	value, exists := c.Get(UserContextKey)
	if !exists {
		return UserContext{}, false
	}

	userCtx, ok := value.(UserContext)

	return userCtx, ok
}

func (m *Middleware) abortWithJSON(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
	return
}
