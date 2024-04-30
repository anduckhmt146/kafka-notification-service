package db

import (
	"log"

	"github.com/anduckhmt146/kakfa-consumer/internal/models"
	"gorm.io/gorm"
)

func autoMigrateSchema(db *gorm.DB) error {
	models := []interface{}{
		&models.User{},
		&models.Notification{},
	}
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			log.Printf("Failed to migrate schema: %v\n", err)
		}
	}
	log.Println("Schema migrated successfully!")
	return nil
}
