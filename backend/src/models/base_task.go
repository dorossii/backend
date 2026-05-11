package models

type TaskTag int

const (
	TaskTagCleaning TaskTag = 0 // 掃除
	TaskTagLaundry  TaskTag = 1 // 洗濯
	TaskTagCooking  TaskTag = 2 // 料理
	TaskTagTrash    TaskTag = 3 // ゴミ捨て
)

type BaseTask struct {
	TaskID          string  `json:"TaskID" gorm:"primaryKey"` // タスクID
	TaskName        string  `json:"TaskName"`                 // タスク名
	DifficultyLevel int     `json:"DifficultyLevel"`          // 難易度
	DueTime         int     `json:"DueTime"`                  // タスク期限
	ImageFlag       bool    `json:"ImageFlag" gorm:"default:false"` // 写真フラグ
	Tags            TaskTag `json:"Tags"`                     // タグ
}