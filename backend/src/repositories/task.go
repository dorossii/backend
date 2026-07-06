package repositories

import (
	"backend/models"
	"time"

	"gorm.io/gorm"
)

type TaskResponse struct {
	TaskID          string    `json:"taskId"`
	UserID          string    `json:"userId"`
	TaskName        string    `json:"taskName"`
	Status          int       `json:"status"`
	Tag             int       `json:"tag"`
	DifficultyLevel int       `json:"level"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"startDate"` // JSONで startDate となるように調整
	EndTime         time.Time `json:"endTime"`
	ImageID         string    `json:"imageId"`
	Message         string    `json:"message"`
}

func GetUserTasks(userID string) ([]TaskResponse, error) {
	var results []TaskResponse

	err := models.DB.Model(&models.Task{}).
		Select(`
            tasks.task_id,  
            tasks.user_id, 
            base_tasks.task_name, 
            tasks.status, 
            base_tasks.tags, 
            base_tasks.difficulty_level,
            base_tasks.description, 
            tasks.start_time as start_date, 
            tasks.end_time, 
            tasks.image_id,
            tasks.message
        `).
		Joins("JOIN base_tasks ON tasks.base_id = base_tasks.base_id").
		Where("tasks.user_id = ?", userID).
		Scan(&results).Error

	return results, err
}

// PendingTaskRow は承認待ちタスク取得のScan用中間構造体
type PendingTaskRow struct {
	TaskID      string
	UserID      string
	TaskName    string
	Tag         int
	Description string
	StartDate   time.Time
	EndTime     time.Time
	ImageID     string
}

// GetPendingTasksForFriends は friendIDs が所有する Status = Pending のタスク一覧を取得する
func GetPendingTasksForFriends(friendIDs []string) ([]PendingTaskRow, error) {
	var results []PendingTaskRow

	if len(friendIDs) == 0 {
		return results, nil
	}

	err := models.DB.Model(&models.Task{}).
		Select(`
            tasks.task_id,
            tasks.user_id,
            base_tasks.task_name,
            base_tasks.tags AS tag,
            base_tasks.description,
            tasks.start_time as start_date,
            tasks.end_time,
            tasks.image_id
        `).
		Joins("JOIN base_tasks ON tasks.base_id = base_tasks.base_id").
		Where("tasks.user_id IN ?", friendIDs).
		Where("tasks.status = ?", models.TaskStatusPending).
		Scan(&results).Error

	return results, err
}

func GetTask(taskID string) (models.Task, error) {
	var task models.Task
	err := models.DB.Model(&models.Task{}).Where("task_id = ?", taskID).First(&task).Error
	return task, err
}

func UpdateTaskImage(taskID string, imageID string) error {
	return models.DB.Model(&models.Task{}).Where("task_id = ?", taskID).Update("image_id", imageID).Error
}

func GetBaseTask(baseID string) (models.BaseTask, error) {
	var basetask models.BaseTask
	err := models.DB.Model(&models.BaseTask{}).Where("base_id = ?", baseID).First(&basetask).Error
	return basetask, err
}

func UpdateTaskStatus(db *gorm.DB, taskID string, status models.TaskStatus) error {
	return db.Model(&models.Task{}).Where("task_id = ?", taskID).Update("status", status).Error
}

func UpdateTaskMessage(tx *gorm.DB, taskID string, message string) error {
	return tx.Model(&models.Task{}).Where("task_id = ?", taskID).Update("message", message).Error
}

func UpdateUserCombo(db *gorm.DB, userID string, combo int, completedAt time.Time) error {
	return db.Model(&models.User{}).Where("user_id = ?", userID).
		Updates(models.User{Combo: combo, LastTaskCompletedAt: &completedAt}).Error
}