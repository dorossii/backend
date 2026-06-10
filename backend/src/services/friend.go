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
	ErrFriendNotFound  = errors.New("指定されたフレンドが存在しません")
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

type FriendInfo struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	HP         int    `json:"hp"`
	DirtLevel  int    `json:"dirtLevel"`
	Icon       string `json:"icon"`
	Background string `json:"background"`
}

// GetFriends は承認済みフレンドの情報一覧を返す
func GetFriends(userID string) ([]FriendInfo, error) {
	// 承認済みフレンドの User レコードを取得
	users, err := repositories.GetFriends(userID)
	if err != nil {
		return nil, err
	}

	// User レコードをレスポンス用の FriendInfo に変換
	friends := make([]FriendInfo, 0, len(users))
	for _, u := range users {
		friends = append(friends, FriendInfo{
			UserID:     u.UserID,
			Name:       u.UserName,
			HP:         u.HealthPoint,
			DirtLevel:  u.DirtLevel,
			Icon:       u.Icon,
			Background: u.BgColor,
		})
	}
	return friends, nil
}

type RescueFriendInfo struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Background string `json:"background"`
	IsRescued  bool   `json:"isrescued"`
}

// GetRescueFriends は承認済みフレンドの一覧をレスキュー状態付きで返す
func GetRescueFriends(userID string) ([]RescueFriendInfo, error) {
	// ① 承認済みフレンドの User レコードを取得
	users, err := repositories.GetFriends(userID)
	if err != nil {
		return nil, err
	}

	// ② 自分が help_targets に登録しているレスキュー対象を取得
	targets, err := repositories.GetHelpTargets(userID)
	if err != nil {
		return nil, err
	}

	// O(1) で検索できるよう FriendID の Set を作成
	rescuedSet := make(map[string]struct{}, len(targets))
	for _, t := range targets {
		rescuedSet[t.FriendID] = struct{}{}
	}

	// フレンド一覧と Set をマージして isrescued を付与
	result := make([]RescueFriendInfo, 0, len(users))
	for _, u := range users {
		// Set に存在すれば true、なければ false
		_, isRescued := rescuedSet[u.UserID]
		result = append(result, RescueFriendInfo{
			UserID:     u.UserID,
			Name:       u.UserName,
			Icon:       u.Icon,
			Background: u.BgColor,
			IsRescued:  isRescued,
		})
	}
	return result, nil
}

var (
	ErrFriendShipNotFound   = errors.New("フレンド関係が存在しません")
	ErrFriendShipNotAccepted = errors.New("承認済みのフレンド関係ではありません")
)

// DeleteFriend はフレンド関係を削除する
func DeleteFriend(userID, friendID string) error {
	// GetFriendShipAny で双方向いずれかのレコードを取得
	fs, err := repositories.GetFriendShipAny(userID, friendID)
	if err != nil {
		return err
	}
	if fs == nil {
		return ErrFriendShipNotFound
	}

	// 承認済みのフレンド関係のみ削除可能
	if fs.Status != models.FriendStatusAccepted {
		return ErrFriendShipNotAccepted
	}

	// 取得したレコードを削除
	return repositories.DeleteFriendShip(fs)
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

// 嫌がらせ設定
func PostAttackerSettings(userID string, targetUser string) error {
	// 空文字ならランダム設定
	if targetUser == "" {
		err := repositories.UpdateAttackerSettings(userID, "")

		if err != nil {
			return err
		}
		return nil
	}

	// 指定ユーザーの場合はフレンドチェック
	friendShip, err := repositories.GetFriendShipAny(userID, targetUser)
	if err != nil {
		return err
	}

	if friendShip == nil {
		return ErrFriendNotFound
	}

	err = repositories.UpdateAttackerSettings(userID, targetUser)
	if err != nil {
		logger.PrintErr("update attacker settings", err)
		return err
	}

	return nil
}
