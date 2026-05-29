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
	// 両方のユーザーが存在するか
	if _, err := repositories.GetUser(userID); err != nil {
		return err
	}
	if _, err := repositories.GetUser(friendID); err != nil {
		return err
	}

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

// フレンドリクエストを承認するエンドポイント
func AcceptFriendRequest(userID, friendID string) error {
	// フレンドシップ取得
	existing, err := repositories.GetFriendShipAny(userID, friendID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("相手からのフレンドリクエストはありません")
	}
	if existing.UserID != friendID {
		return errors.New("相手からのフレンドリクエストを承認する権限がありません")
	}

	// 承認
	existing.Status = models.FriendStatusAccepted

	// 更新
	return repositories.UpdateFriendShip(existing)
}

type FriendRequest struct {
	RequestUser string `json:"RequestUser"`
	CreatedAt   int64  `json:"CreatedAt"`
}

// フレンド申請一覧取得
func GetFriendRequests(userID string) ([]FriendRequest, error) {
	// 相手からのフレンドリクエスト
	FriendReqests := []FriendRequest{}

	// 相手からのフレンドリクエスト
	getReqs, err := repositories.GetIncomingFriendShipsByStatus(userID, models.FriendStatusPending)

	// エラー
	if err != nil {
		return FriendReqests, err
	}

	// 相手からのフレンドリクエストを取得
	for _, req := range getReqs {
		FriendReqests = append(FriendReqests, FriendRequest{
			RequestUser: req.UserID,
			CreatedAt:   req.CreatedAt,
		})
	}

	return FriendReqests, nil
}
