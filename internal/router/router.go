package router

import (
	"github.com/AlexShmak/wb_test_task_l0/internal/auth"
	"github.com/AlexShmak/wb_test_task_l0/internal/handlers"
	"log/slog"
	"time"

	"github.com/AlexShmak/wb_test_task_l0/internal/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(storage *storage.PostgresStorage, logger *slog.Logger, jwtService *auth.JWTService) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(gin.Recovery())

	handler := handlers.NewHandler(storage, logger, jwtService)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handler.RegisterUserHandler)
		authGroup.POST("/login", handler.LoginHandler)
		authGroup.POST("/logout", handler.LogoutHandler)
	}

	api := router.Group("/api")
	api.Use(handler.AuthMiddleware())
	{
		api.GET("/orders/:id", handler.GetOrderByIDHandler)
	}

	return router
}
