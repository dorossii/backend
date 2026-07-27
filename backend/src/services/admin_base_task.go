package services

import (
	"backend/models"
	"backend/repositories"
	"errors"

	"gorm.io/gorm"
)

var ErrBaseTaskNotFound = errors.New("指定されたBaseTaskが存在しません")

func AdminCreateBaseTask(bt *models.BaseTask) error {
	return repositories.CreateBaseTask(bt)
}

func AdminUpdateBaseTask(baseID string, updates *models.BaseTask) (*models.BaseTask, error) {
	existing, err := repositories.GetBaseTask(baseID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBaseTaskNotFound
	}
	if err != nil {
		return nil, err
	}

	updates.BaseID = existing.BaseID
	if err := repositories.UpdateBaseTask(updates); err != nil {
		return nil, err
	}

	return updates, nil
}

func AdminDeleteBaseTask(baseID string) error {
	existing, err := repositories.GetBaseTask(baseID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrBaseTaskNotFound
	}
	if err != nil {
		return err
	}

	return repositories.DeleteBaseTask(existing.BaseID)
}

func AdminListBaseTasks() ([]models.BaseTask, error) {
	return repositories.ListBaseTasks()
}
