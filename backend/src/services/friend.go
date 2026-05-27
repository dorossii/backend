package services

import (
	"backend/logger"
	"backend/models"
	"backend/repositories"
	"errors"
)

var (
	ErrAlreadySent     = errors.New("すでにフレンド申請済みです")
	ErrAlreadyReceived = errors.New("相手からすでに申請が届いています")
	ErrFriendNotFound = errors.New("指定されたフレンドが存在しません")
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

func PostAttackerSettings(userID string, targetUser string) error {
	// 空文字ならランダム設定
	if targetUser == "" {
		err := repositories.UpdateAttackerSettings(userID, "")

		if err != nil {
			return err
		}
		return nil
	}

	// 指定ユーザーの場合は friend check
	isFriend, err := repositories.IsFriend(userID, targetUser)
	if err != nil {
		return err
	}

	if !isFriend {
		return ErrFriendNotFound
	}

	err = repositories.UpdateAttackerSettings(userID, targetUser)
	if err != nil {
		logger.PrintErr("update attacker settings", err)
		return err
	}

	return nil
}