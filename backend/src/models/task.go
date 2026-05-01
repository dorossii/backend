package models

import "time"

type TaskStatus int

const (
	TaskStatusImcomplete TaskStatus = 0 // 未完了
	TaskStatusPending    TaskStatus = 1 // 承認待ち
	TaskStatusCompleted  TaskStatus = 2 // 完了
)

type Task struct {
	TaskID       string    `json:"TaskID"; gorm:"primaryKey"`     // タスクのID
	UserID       string    `json:"UserID"`                        // タスクを持つユーザのID
	Status       int       `json:"TaskName"`                      // タスクの名前
	StartTime    time.Time `json:"StartTime"`                     // タスクが付与された時間
	EndTime      time.Time `json:"EndTime"`                       // タスクが終了する時間
	ImageID      string    `json:"ImageID"`                       // 画像のID (タスクの内容を表す画像)
	RequireImage bool      `json:"RequireImage"; default:"false"` // 画像の要求の有無
}
