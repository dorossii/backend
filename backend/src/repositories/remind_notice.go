package repositories

import (
	"backend/models"
)

func CreateRemindNotiec(notice *models.RemindNotice) error {
	return models.DB.Create(notice).Error
}