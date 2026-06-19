package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	return models.DB.Create(user).Error
}

func GetUser(UserID string) (*models.User, error) {
	var user models.User
	err := models.DB.First(&user, "user_id = ?", UserID).Error
	return &user, err
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

func UpdateAttackerSettingsTx(tx *gorm.DB, userID string, targetUser string) error {
	err := tx.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("target_user", targetUser).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAttackerSettings は userID の target_user を返す	
func GetAttackerSettings(userID string) (string, error) {
	var user models.User
	err := models.DB.Select("target_user").First(&user, "user_id = ?", userID).Error
	if err != nil {
		return "", err
	}
	
	return user.TargetUser, nil
}

