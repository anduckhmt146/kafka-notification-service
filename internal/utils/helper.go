package utils

import (
	"github.com/anduckhmt146/kakfa-consumer/internal/constants"
	"github.com/anduckhmt146/kakfa-consumer/internal/models"
)

func FindUserByID(id int, users []models.User) (models.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, constants.ErrUserNotFoundInProducer
}
