package services

import (
	"backend/repositories"
	"fmt"
	"sort"
)

// NoticeResponse はアプリ内通知の共通表現
type NoticeResponse struct {
	SenderType string `json:"senderType"` // "official" または送信者のユーザーID
	Title      string `json:"title"`
	NotifiedAt int64  `json:"notifiedAt"`
}

// userDisplayName はユーザーIDから表示名を解決する。取得できない場合は userID をそのまま返す
func userDisplayName(userID string) string {
	user, err := repositories.GetUser(userID)
	if err != nil {
		return userID
	}
	return user.UserName
}

// GetNotices は指定ユーザー宛の通知一覧を通知時刻の降順で返す
// RemindNotice は対象外（別エンドポイントで扱う想定）
func GetNotices(userID string) ([]NoticeResponse, error) {
	notices := []NoticeResponse{}

	trashNotices, err := repositories.GetTrashNoticesByReceiver(userID)
	if err != nil {
		return nil, err
	}
	for _, n := range trashNotices {
		notices = append(notices, NoticeResponse{
			SenderType: n.SenderID,
			Title:      fmt.Sprintf("%sさんから汚さ%d分の攻撃が届きました", userDisplayName(n.SenderID), n.Count),
			NotifiedAt: n.CreatedAt.Unix(),
		})
	}

	rescueNotices, err := repositories.GetRescueNoticesByHelper(userID)
	if err != nil {
		return nil, err
	}
	for _, n := range rescueNotices {
		notices = append(notices, NoticeResponse{
			SenderType: n.TargetID,
			Title:      fmt.Sprintf("%sさんをレスキューしました", userDisplayName(n.TargetID)),
			NotifiedAt: n.CreatedAt.Unix(),
		})
	}

	helpedNotices, err := repositories.GetHelpedNoticesByTarget(userID)
	if err != nil {
		return nil, err
	}
	for _, n := range helpedNotices {
		notices = append(notices, NoticeResponse{
			SenderType: n.HelperID,
			Title:      fmt.Sprintf("%sさんにレスキューされました", userDisplayName(n.HelperID)),
			NotifiedAt: n.CreatedAt.Unix(),
		})
	}

	sort.Slice(notices, func(i, j int) bool {
		return notices[i].NotifiedAt > notices[j].NotifiedAt
	})

	return notices, nil
}
