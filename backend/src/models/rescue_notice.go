package models

import "time"

type RescueNotice struct {
	TargetID string    `json:"TargetID" gorm:"primaryKey"` // タスクしてない側ID
	HelperID string    `json:"HelperID" gorm:"primaryKey"` // 助ける側ID
	IsRead   bool      `json:"IsRead" gorm:"default:false"` // 既読状態
	CreatedAt time.Time `json:"CreatedAt"` // 作成時間
}