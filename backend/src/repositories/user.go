package repositories

import (
	"backend/models"
)

func CreateUser(user *models.User) error {
	return models.DB.Create(user).Error
}

func GetUser(UserID string) (*models.User, error) {
	var user models.User
	err := models.DB.First(&user, "user_id = ?", UserID).Error
	return &user, err
}
