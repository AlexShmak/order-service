package router

import (
	"github.com/AlexShmak/order-service/internal/config"
	"github.com/AlexShmak/order-service/internal/kafka"
	"github.com/AlexShmak/order-service/internal/storage/cache"
	"log/slog"
	"time"

	"github.com/AlexShmak/order-service/internal/auth"
	"github.com/AlexShmak/order-service/internal/handlers"

	"github.com/AlexShmak/order-service/internal/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(storage *storage.PostgresStorage, logger *slog.Logger, jwtService *auth.JWTService, cfg *config.Config, producer *kafka.Producer, redisCache *cache.RedisStorage) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(gin.Recovery())

	handler := handlers.NewHandler(storage, logger, jwtService, cfg, producer, redisCache)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handler.RegisterHandler)
		authGroup.POST("/login", handler.LoginHandler)
		authGroup.POST("/logout", handler.LogoutHandler)
	}

	api := router.Group("/api")
	api.Use(handler.AuthMiddleware())
	{
		api.GET("/orders/:id", handler.GetOrderByIDHandler)
		api.POST("/orders", handler.CreateOrderHandler)
	}

	return router
}
