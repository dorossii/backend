package services

import (
	"backend/models"
	"backend/repositories"
	"errors"
)

var (
	ErrAlreadySent     = errors.New("すでにフレンド申請済みです")
	ErrAlreadyReceived = errors.New("相手からすでに申請が届いています")
)

func SendFriendRequest(userID, friendID string) error {
	// 自分 → 相手のレコードを確認
	existing, err := repositories.GetFriendShip(userID, friendID)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrAlreadySent
	}

	// 相手 → 自分のレコードを確認
	reverse, err := repositories.GetFriendShip(friendID, userID)
	if err != nil {
		return err
	}
	if reverse != nil {
		return ErrAlreadyReceived
	}

	return repositories.CreateFriendShip(&models.FriendShips{
		UserID:   userID,
		FriendID: friendID,
		Status:   models.FriendStatusPending,
	})
}
