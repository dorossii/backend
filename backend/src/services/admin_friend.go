package services

import (
	"backend/models"
	"backend/repositories"
)

// AdminForceFriendShip は既存関係の有無に関わらずフレンド関係を Accepted 状態にする
func AdminForceFriendShip(userID, friendID string) error {
	existing, err := repositories.GetFriendShipAny(userID, friendID)
	if err != nil {
		return err
	}

	if existing == nil {
		return repositories.CreateFriendShip(&models.FriendShips{
			UserID:   userID,
			FriendID: friendID,
			Status:   models.FriendStatusAccepted,
		})
	}

	existing.Status = models.FriendStatusAccepted
	return repositories.UpdateFriendShip(existing)
}

func AdminUpdateFriendShipStatus(userID, friendID string, status models.FriendStatus) error {
	existing, err := repositories.GetFriendShipAny(userID, friendID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrFriendShipNotFound
	}

	existing.Status = status
	return repositories.UpdateFriendShip(existing)
}

func AdminDeleteFriendShip(userID, friendID string) error {
	existing, err := repositories.GetFriendShipAny(userID, friendID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrFriendShipNotFound
	}

	return repositories.DeleteFriendShip(existing)
}
