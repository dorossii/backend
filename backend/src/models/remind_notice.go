package models

import "time"

type SenderType string

const (
	SenderTypeOfficial SenderType = "official" // 公式通知
	SenderTypeFriend   SenderType = "friend"   // フレンド通知
)

type RemindNotice struct {
	UserID      string     `json:"UserID" gorm:"primaryKey"`           // ユーザーID
	SenderType  SenderType `json:"SenderType"`                         // 通知の送信元
	Title       string     `json:"Title"`                              // タイトル
	NotifiedAt  time.Time  `json:"NotifiedAt"`                         // 受信時間
	IsRead      bool       `json:"IsRead" gorm:"default:false"`        // 既読状態
}