package repositories

import (
	"backend/models"
)

// レスキューを更新する
func UpdateRescuerSettings(userID, targetUser string) error {
	return models.DB.Model(&models.HelpTargets{}).Create(&models.HelpTargets{UserID: userID, FriendID: targetUser}).Error
}

// レスキューを削除する
func DeleteRescuerSettings(userID string) error {
	return models.DB.Where("user_id = ? ", userID).Delete(&models.HelpTargets{}).Error
}
