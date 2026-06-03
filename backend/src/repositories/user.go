package repositories

import (
	"backend/models"
)

func CreateUser(user *models.User) error {
	return models.DB.Create(user).Error
}

func GetUser(userID string) (models.User, error) {
	var user models.User
	err := models.DB.Model(&models.User{}).Where("user_id = ?", userID).Find(&user).Error
	return user, err
}

func UpdateDirtLevel(userID string, diff int) error {
	var user models.User

	err := models.DB.Model(&models.User{}).First(&user, "user_id = ?", userID).Error
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

	return models.DB.Model(&models.User{}).Where("user_id = ?", userID).Update("dirt_level", newDirt).Error
}

