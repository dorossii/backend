package models

import "time"

type TrashNotice struct {
	SenderID   string    `json:"SenderID" gorm:"primaryKey"` // ゴミを投げた人
	ReceiverID string    `json:"ReceiverID" gorm:"primaryKey"` // ゴミを投げられた人
	Count      int       `json:"Count" gorm:"default:0"` // ゴミを投げられた数
	CreatedAt  time.Time `json:"CreatedAt"` // 作成時間
}