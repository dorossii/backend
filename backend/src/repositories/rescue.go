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

// レスキュー設定をトランザクションで更新する
func PostRescuerSettingsTx(userID string, targetUsers []string) error {
	return models.DB.Transaction(func(tx *gorm.DB) error {
		// 既存のレスキュー設定を削除
		if err := tx.Where("user_id = ?", userID).Delete(&models.HelpTargets{}).Error; err != nil {
			return err
		}

		// 空ならランダム設定（FriendID空のレコード1件）
		if len(targetUsers) == 0 {
			return tx.Create(&models.HelpTargets{UserID: userID, FriendID: ""}).Error
		}

		// 指定ユーザーを一括インサート
		targets := make([]models.HelpTargets, 0, len(targetUsers))
		for _, targetUser := range targetUsers {
			targets = append(targets, models.HelpTargets{UserID: userID, FriendID: targetUser})
		}
		return tx.Create(&targets).Error
	})
}
