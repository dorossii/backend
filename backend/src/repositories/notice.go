package repositories

import (
	"backend/models"
)

func CreateRemindNotiec(notice *models.RemindNotice) error {
	return models.DB.Create(notice).Error
}

func CreateTrashNotice(notice *models.TrashNotice) error {
	return models.DB.Create(notice).Error
}