package models

import "time"

type SenderType string

type RemindNotice struct {
	NoticeID   string    `json:"NoticeID" gorm:"primaryKey"`  // 通知ID
	UserID     string    `json:"UserID"`                      // ユーザーID
	SenderID   string    `json:"SenderID"`                    // 通知の送信元
	Title      string    `json:"Title"`                       // タイトル
	NotifiedAt time.Time `json:"NotifiedAt"`                  // 受信時間
	IsRead     bool      `json:"IsRead" gorm:"default:false"` // 既読状態
}
