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
	existing, err := repositories.GetFriendShipAny(userID, friendID)
	if err != nil {
		return err
	}
	if existing != nil {
		if existing.UserID == userID {
			return ErrAlreadySent
		}
		return ErrAlreadyReceived
	}

	return repositories.CreateFriendShip(&models.FriendShips{
		UserID:   userID,
		FriendID: friendID,
		Status:   models.FriendStatusPending,
	})
}
