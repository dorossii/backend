package repositories

import (
	"backend/models"
	"time"
)

type TaskResponse struct {
    TaskID      string    `json:"taskId"`
    UserID      string    `json:"userId"`
    TaskName    string    `json:"taskName"`
    Status      int       `json:"status"`
    Tag         int        `json:"tag"`
    Description string    `json:"description"`
    StartDate   time.Time `json:"startDate"` // JSONで startDate となるように調整
    EndTime     time.Time `json:"endTime"`
    ImageID     string    `json:"imageId"`
}

func GetTasks(userID string) ([]models.Task, error) {
	//taskとbaseを結合して、ユーザーIDに紐づくタスクを取得する
	var tasks []models.Task
	if err := models.DB.Preload("Base").Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
