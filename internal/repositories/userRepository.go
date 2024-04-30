package repositories

import (
	"github.com/anduckhmt146/kakfa-consumer/internal/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByID(id int) (*models.User, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (repo *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	result := repo.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
