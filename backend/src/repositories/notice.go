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