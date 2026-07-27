package services

import (
	"backend/models"
	"backend/repositories"
	"errors"

	"gorm.io/gorm"
)

func AdminListTasks(userID string) ([]models.Task, error) {
	return repositories.ListTasks(userID)
}

func AdminCreateTask(task *models.Task) error {
	return repositories.CreateTask(task)
}

func AdminUpdateTask(taskID string, updates *models.Task) (*models.Task, error) {
	existing, err := repositories.GetTask(taskID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}

	updates.TaskID = existing.TaskID
	if err := repositories.UpdateTask(updates); err != nil {
		return nil, err
	}

	return updates, nil
}

func AdminUpdateTaskStatus(taskID string, status models.TaskStatus) (*models.Task, error) {
	existing, err := repositories.GetTask(taskID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := repositories.UpdateTaskStatus(models.DB, existing.TaskID, status); err != nil {
		return nil, err
	}
	existing.Status = status

	return &existing, nil
}

func AdminDeleteTask(taskID string) error {
	existing, err := repositories.GetTask(taskID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrTaskNotFound
	}
	if err != nil {
		return err
	}

	return repositories.DeleteTask(existing.TaskID)
}

// AdminGetTaskImagePath は指定Taskに紐づく画像の実ファイルパスを返す
func AdminGetTaskImagePath(taskID string) (filePath string, err error) {
	task, err := repositories.GetTask(taskID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", ErrTaskNotFound
	}
	if err != nil {
		return "", err
	}
	if task.ImageID == "" {
		return "", ErrTaskNotFound
	}

	filePath, _, err = GetTaskImage(task.ImageID)
	return filePath, err
}
