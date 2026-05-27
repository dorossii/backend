package repositories

import (
	"backend/models"
)

func CreateUser(user *models.User) error {
	return models.DB.Create(user).Error
}

func UpdateAttackerSettings(userID string, targetUser string) error {
	err := models.DB.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("target_user", targetUser).Error
	if err != nil {
		return err
	}

	return nil
}
