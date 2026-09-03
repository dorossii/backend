package repositories

import (
	"backend/models"
	"time"

	"gorm.io/gorm"
)

type TaskResponse struct {
	TaskID          string   `json:"taskId"`
	UserID          string   `json:"userId"`
	TaskName        string   `json:"taskName"`
	Status          int      `json:"status"`
	Tag             int      `json:"tag"`
	DifficultyLevel int      `json:"level"`
	Description     string   `json:"description"`
	StartDate       int64    `json:"startDate"`
	EndTime         int64    `json:"endTime"`
	ImageID         string   `json:"imageId"`
	Message         string   `json:"message"`
}

func GetUserTasks(userID string) ([]TaskResponse, error) {
	results := []TaskResponse{}

	err := models.DB.Model(&models.Task{}).
		Select(`
            tasks.task_id,  
            tasks.user_id, 
            base_tasks.task_name, 
            tasks.status, 
            base_tasks.tags AS tag,
            base_tasks.difficulty_level,
            base_tasks.description, 
            CAST(UNIX_TIMESTAMP(tasks.start_time) AS SIGNED) as start_date, 
            CAST(UNIX_TIMESTAMP(tasks.end_time) AS SIGNED) as end_time, 
            base_tasks.image_id,
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

func CreateTask(task *models.Task) error {
	return models.DB.Create(task).Error
}

func UpdateTask(task *models.Task) error {
	return models.DB.Save(task).Error
}

func DeleteTask(taskID string) error {
	return models.DB.Where("task_id = ?", taskID).Delete(&models.Task{}).Error
}

// ListTasks は userID が空文字なら全件、指定されていればそのユーザーのタスク一覧を返す
func ListTasks(userID string) ([]models.Task, error) {
	var tasks []models.Task
	query := models.DB.Model(&models.Task{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Find(&tasks).Error
	return tasks, err
}

// AdminTaskRow は管理画面のタスク一覧向けに Task と BaseTask の内容を突き合わせた行
type AdminTaskRow struct {
	TaskID          string           `json:"TaskID"`
	BaseID          string           `json:"BaseID"`
	UserID          string           `json:"UserID"`
	Status          models.TaskStatus `json:"Status"`
	StartTime       time.Time        `json:"StartTime"`
	EndTime         time.Time        `json:"EndTime"`
	ImageID         string           `json:"ImageID"`
	RequireImage    bool             `json:"RequireImage"`
	Message         string           `json:"Message"`
	TaskName        string           `json:"TaskName"`
	Description     string           `json:"Description"`
	DifficultyLevel int              `json:"DifficultyLevel"`
	DueTime         int              `json:"DueTime"`
	Tags            models.TaskTag   `json:"Tags"`
}

// AdminListTasksWithBaseTask は userID が空文字なら全件、指定されていればそのユーザーのタスク一覧を
// BaseTask の内容(タスク名・タグなど)と突き合わせて返す
func AdminListTasksWithBaseTask(userID string) ([]AdminTaskRow, error) {
	var rows []AdminTaskRow

	query := models.DB.Model(&models.Task{}).
		Select(`
            tasks.task_id,
            tasks.base_id,
            tasks.user_id,
            tasks.status,
            tasks.start_time,
            tasks.end_time,
            tasks.image_id,
            tasks.require_image,
            tasks.message,
            base_tasks.task_name,
            base_tasks.description,
            base_tasks.difficulty_level,
            base_tasks.due_time,
            base_tasks.tags
        `).
		Joins("JOIN base_tasks ON tasks.base_id = base_tasks.base_id")

	if userID != "" {
		query = query.Where("tasks.user_id = ?", userID)
	}

	err := query.Scan(&rows).Error
	return rows, err
}

// フレンド且つRemindNoticeに保存されていないユーザーをランダムで一件取得する
// 対象が存在しない場合は空文字列, nil を返す（エラーではない）
func GetMessageUserId(userID string) (string, error) {
	var friendID string

	err := models.DB.
		Raw(`
			SELECT CASE WHEN friend_ships.user_id = ? THEN friend_ships.friend_id ELSE friend_ships.user_id END AS friend_id
			FROM friend_ships
			WHERE (friend_ships.user_id = ? OR friend_ships.friend_id = ?)
			AND friend_ships.status = ?
			AND NOT EXISTS (
				SELECT 1
				FROM remind_notices
				WHERE remind_notices.user_id = CASE WHEN friend_ships.user_id = ? THEN friend_ships.friend_id ELSE friend_ships.user_id END
				AND remind_notices.sender_id = ?
			)
			ORDER BY RAND()
			LIMIT 1
		`, userID, userID, userID, models.FriendStatusAccepted, userID, userID).
		Scan(&friendID).Error

	return friendID, err
}

// 複数まとめてステータス更新
func UpdateTasksStatus(tx *gorm.DB, taskIDs []string, status models.TaskStatus) error {
	return tx.Model(&models.Task{}).Where("task_id IN ?", taskIDs).Update("status", status).Error
}
