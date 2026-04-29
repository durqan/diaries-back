package main

import (
	"diary-backend/internal/config"
	"diary-backend/internal/models"
	"diary-backend/internal/routes"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	if err := db.AutoMigrate(&models.Entry{}); err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	r := routes.SetupRouter(db)
	r.Run(":" + cfg.Port)
}
