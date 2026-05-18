package repositories

import (
	"backend/models"
)

func CreateUser(user *models.User) error {
	return models.DB.Create(user).Error
}
