package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	return models.DB.Create(user).Error
}

func CreateUserTx(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

func CreateUserRoomTx(tx *gorm.DB, userRoom *models.UserRoom) error {
	return tx.Create(userRoom).Error
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

func UpdateDirtLevel(tx *gorm.DB, userID string, diff int) error {
	var user models.User

	err := tx.Model(&models.User{}).First(&user, "user_id = ?", userID).Error
	if err != nil {
		return err
	}

	// 計算後の値が0以下700以上ならぱわーでねじ曲げる
	newDirt := user.DirtLevel + diff

	if newDirt < 0 {
		newDirt = 0
	}

	if newDirt > 700 {
		newDirt = 700
	}

	return tx.Model(&models.User{}).Where("user_id = ?", userID).Update("dirt_level", newDirt).Error
}

