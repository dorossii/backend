package repositories

import (
	"backend/models"
)

func CreateBaseTask(bt *models.BaseTask) error {
	return models.DB.Create(bt).Error
}

func UpdateBaseTask(bt *models.BaseTask) error {
	return models.DB.Save(bt).Error
}

func DeleteBaseTask(baseID string) error {
	return models.DB.Where("base_id = ?", baseID).Delete(&models.BaseTask{}).Error
}

func ListBaseTasks() ([]models.BaseTask, error) {
	var baseTasks []models.BaseTask
	err := models.DB.Find(&baseTasks).Error
	return baseTasks, err
}
