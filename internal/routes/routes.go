package routes

import (
	"diary-backend/internal/handlers"
	"diary-backend/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	limiterGin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://80.64.27.84", "http://80.64.27.84:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	rate := limiter.Rate{Period: 1 * time.Second, Limit: 30}
	store := memory.NewStore()
	rateLimiter := limiterGin.NewMiddleware(limiter.New(store, rate))
	r.Use(rateLimiter)

	authHandler := handlers.NewAuthHandler(db)
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	entryHandler := handlers.NewEntryHandler(db)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/entries", entryHandler.GetEntries)
		protected.POST("/entries", entryHandler.CreateEntry)
		protected.GET("/entries/:id", entryHandler.GetEntry)
		protected.PUT("/entries/:id", entryHandler.UpdateEntry)
		protected.DELETE("/entries/:id", entryHandler.DeleteEntry)
		protected.GET("/me", authHandler.GetMe)
	}

	return r
}
