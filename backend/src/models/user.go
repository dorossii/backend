package models

import "time"

type User struct {
	UserID      string    `json:"UserID" gorm:"primaryKey"`  // ユーザのID
	UserName    string    `json:"UserName"`                  // ユーザの名前
	BirthDate   time.Time `json:"BirthDate"`                 // ユーザの誕生日
	Mailadress  string    `json:"Mailadress"`                // ユーザのメールアドレス
	HealthPoint int       `json:"HealthPoint" default:"100"` // ユーザの体力
	DirtLevel   int       `json:"DirtLevel" default:"0"`     // ユーザの汚れレベル
	Combo       int       `json:"Combo" default:"0"`         // 連続タスク達成数
	TargetUser  string    `json:"TargetUser"`                // 嫌がらせ対象のユーザのUserID
	BgColor     string    `json:"BgColor" default:"#ffb6c1"` // ユーザの背景色
}
