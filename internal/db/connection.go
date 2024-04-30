package db

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	dbHost := viper.GetString("mysql.DB_HOST")
	dbPort := viper.GetString("mysql.DB_PORT")
	dbUser := viper.GetString("mysql.DB_USER")
	dbPassword := viper.GetString("mysql.DB_PASSWORD")
	dbName := viper.GetString("mysql.DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Println("Failed to connect to MySQL database!")
		return nil, err
	}

	err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET = `utf8mb4` COLLATE = `utf8mb4_general_ci`;", dbName)).Error
	if err != nil {
		log.Println("Failed to create database!")
		return nil, err
	}

	err = db.Exec(fmt.Sprintf("USE %s;", dbName)).Error
	if err != nil {
		log.Println("Failed to use database!")
		return nil, err
	}

	autoMigrateSchema(db)
	return db, nil
}
