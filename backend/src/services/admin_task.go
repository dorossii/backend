package services

import (
	"backend/models"
	"backend/repositories"
	"backend/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

func AdminListTasks(userID string) ([]repositories.AdminTaskRow, error) {
	return repositories.AdminListTasksWithBaseTask(userID)
}

func AdminCreateTask(task *models.Task) error {
	baseTask, err := repositories.GetBaseTask(task.BaseID)
	if err != nil {
		return err
	}

	if task.TaskID == "" {
		uuid, err := utils.Genid()
		if err != nil {
			return err
		}
		task.TaskID = uuid
	}

	// StartTime/EndTime が未指定の場合は BaseTask の期限(日数)から自動計算する
	if task.StartTime.IsZero() {
		task.StartTime = utils.NowJST()
	}
	if task.EndTime.IsZero() {
		task.EndTime = task.StartTime.Add(time.Duration(baseTask.DueTime) * 24 * time.Hour)
	}

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
