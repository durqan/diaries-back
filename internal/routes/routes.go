package routes

import (
	"diary-backend/internal/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	limiterGin "github.com/ulule/limiter/v3/drivers/middleware/gin" // 👈 этот импорт
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://80.64.27.84", "http://80.64.27.84:5173", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  10,
	}
	store := memory.NewStore()

	rateLimiterMiddleware := limiterGin.NewMiddleware(limiter.New(store, rate))
	r.Use(rateLimiterMiddleware)

	handler := handlers.NewEntryHandler(db)

	api := r.Group("/api")
	{
		api.GET("/entries", handler.GetEntries)
		api.POST("/entries", handler.CreateEntry)
		api.GET("/entries/:id", handler.GetEntry)
		api.PUT("/entries/:id", handler.UpdateEntry)
		api.DELETE("/entries/:id", handler.DeleteEntry)
	}

	return r
}
