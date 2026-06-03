package models

type HelpTargets struct {
	UserID   string `json:"UserID" gorm:"primaryKey"`   // ユーザのID
	FriendID string `json:"FriendID" gorm:"primaryKey"` // 嫌がらせ対象のユーザのID
}
