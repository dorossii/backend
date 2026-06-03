package models

import "time"

type TrashNotice struct {
	NoticeID   string    `json:"NoticeID" gorm:"primaryKey"` // 通知ID
	SenderID   string    `json:"SenderID"`                   // ゴミを投げた人
	ReceiverID string    `json:"ReceiverID"`                 // ゴミを投げられた人
	Count      int       `json:"Count" gorm:"default:0"`     // ゴミを投げられた数
	CreatedAt  time.Time `json:"CreatedAt"`                  // 作成時間
}
