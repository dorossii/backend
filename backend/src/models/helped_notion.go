package models

import "time"

type HelpedNotice struct {
	NoticeID  string    `json:"NoticeID" gorm:"primaryKey"`  // 通知ID
	TargetID  string    `json:"TargetID"`                    // タスクしてない側ID
	HelperID  string    `json:"HelperID"`                    // 助ける側ID
	IsRead    bool      `json:"IsRead" gorm:"default:false"` // 既読状態
	CreatedAt time.Time `json:"CreatedAt"`                   // 作成時間
}
