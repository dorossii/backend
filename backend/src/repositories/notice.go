package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

func CreateRemindNotiec(notice *models.RemindNotice) error {
	return models.DB.Create(notice).Error
}

func CreateTrashNotice(tx *gorm.DB, notice *models.TrashNotice) error {
	return tx.Create(notice).Error
}

// GetTrashNoticesByReceiver は指定ユーザーが受け取った TrashNotice 一覧を返す
func GetTrashNoticesByReceiver(userID string) ([]models.TrashNotice, error) {
	var notices []models.TrashNotice
	err := models.DB.Where("receiver_id = ?", userID).Find(&notices).Error
	return notices, err
}

// GetRescueNoticesByHelper は指定ユーザーがレスキューした側として作成された RescueNotice 一覧を返す
func GetRescueNoticesByHelper(userID string) ([]models.RescueNotice, error) {
	var notices []models.RescueNotice
	err := models.DB.Where("helper_id = ?", userID).Find(&notices).Error
	return notices, err
}

// GetHelpedNoticesByTarget は指定ユーザーが助けられた側として作成された HelpedNotice 一覧を返す
func GetHelpedNoticesByTarget(userID string) ([]models.HelpedNotice, error) {
	var notices []models.HelpedNotice
	err := models.DB.Where("target_id = ?", userID).Find(&notices).Error
	return notices, err
}