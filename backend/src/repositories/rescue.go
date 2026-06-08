package repositories

import (
	"backend/models"
)

// UpdateRescuerSettings はレスキューセッティングを更新する
func UpdateRescuerSettings(userID, targetUser string) error {
	return models.DB.Model(&models.User{}).Where("UserID = ?", userID).Update("FriendID", targetUser).Error
}
