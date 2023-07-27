package db

import (
	"fmt"
	"service1/pkg/entity"
	"service1/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() (*gorm.DB, error) {
	config, err := utils.LoadConfig("./")
	fmt.Println("DSN value:", config.DSN)
	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	DB = db
	DB.AutoMigrate(&entity.OtpKey{}, &entity.Signup{}, &entity.Admin{}, &entity.User{}, &entity.Address{})
	return db, nil
}
