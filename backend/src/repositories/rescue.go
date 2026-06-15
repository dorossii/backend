package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

// レスキューを更新する
func UpdateRescuerSettings(tx *gorm.DB, helpTargets []models.HelpTargets) error {
    
    // スライスを渡すことで、GORMは自動的に Bulk Insert クエリを生成します
    return tx.Model(&models.HelpTargets{}).Create(&helpTargets).Error
}

// レスキューを削除する
func DeleteRescuerSettings(tx *gorm.DB, userID string) error {
	return tx.Where("user_id = ? ", userID).Delete(&models.HelpTargets{}).Error
}


