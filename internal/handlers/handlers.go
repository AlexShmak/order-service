package handlers

import (
	"errors"
	"github.com/AlexShmak/wb_test_task_l0/internal/auth"
	"github.com/AlexShmak/wb_test_task_l0/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
)

type Handler struct {
	Storage    *storage.PostgresStorage
	Logger     *slog.Logger
	JWTService *auth.JWTService
}

func NewHandler(s *storage.PostgresStorage, l *slog.Logger, jwt *auth.JWTService) *Handler {
	return &Handler{
		Storage:    s,
		Logger:     l,
		JWTService: jwt,
	}
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			h.Logger.Error("no access token found in cookies", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		token, err := h.JWTService.ValidateAccessToken(tokenString)
		if err != nil || !token.Valid {
			h.Logger.Error("invalid token", slog.String("error", err.Error()))
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userId, err := h.JWTService.GetUserIdFromToken(token)
		if err != nil {
			h.Logger.Error("could not get user ID from token", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}

func (h *Handler) RegisterUserHandler(c *gin.Context) {
	// Implementation for user registration
}
func (h *Handler) LoginHandler(c *gin.Context) {
	// Implementation for user login
}
func (h *Handler) LogoutHandler(c *gin.Context) {
	// Implementation for user logout
}
func (h *Handler) GetOrderByIDHandler(c *gin.Context) {
	// Implementation for getting order by ID
}
