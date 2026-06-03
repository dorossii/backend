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
    Tag         int       `json:"tag"`
    DifficultyLevel int   `json:"level"`
    Description string    `json:"description"`
    StartDate   time.Time `json:"startDate"` // JSONで startDate となるように調整
    EndTime     time.Time `json:"endTime"`
    ImageID     string    `json:"imageId"`
    Message     string    `json:"message"`
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
