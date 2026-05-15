package models

type FriendStatus int

const (
	FriendStatusPending  FriendStatus = 0 // 申請中
	FriendStatusAccepted FriendStatus = 1 // 友達 (成立)
	FriendStatusRejected FriendStatus = 2 // 友達 (拒否)
)

type FriendShips struct {
	UserID   string       `json:"UserID" gorm:"primaryKey"`   // ユーザのID
	FriendID string       `json:"FriendID" gorm:"primaryKey"` // 友達のユーザID
	Status   FriendStatus `json:"Status"`                      // 友達関係の状態
}
