package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

// レスキューを更新する
func UpdateRescuerSettings(tx *gorm.DB, userID, targetUser string) error {
	return tx.Model(&models.HelpTargets{}).Create(&models.HelpTargets{UserID: userID, FriendID: targetUser}).Error
}

// レスキューを削除する
func DeleteRescuerSettings(tx *gorm.DB, userID string) error {
	return tx.Where("user_id = ? ", userID).Delete(&models.HelpTargets{}).Error
}
