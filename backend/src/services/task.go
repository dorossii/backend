package services

import (
	"backend/logger"
	"backend/models"
	"backend/repositories"
	"time"
	"errors"

	"github.com/google/uuid"
)

type TaskService struct{}

// errors
var ErrFriendNotFound = errors.New("指定されたフレンドが存在しません")

// 煽りメッセージの登録
func (s *TaskService) PostTaskTauntMessage(userId string, friendId string, msg string) error {
	// フレンド存在確認
	isFriend, err := repositories.IsFriend(userId, friendId)
	if err != nil {
		return err
	}

	if !isFriend {
		return ErrFriendNotFound
	}

	// メッセージの登録
	notice := &models.RemindNotice{
		NoticeID:   uuid.NewString(),
		UserID:     userId,
		SenderID:   friendId,
		Title:      msg,
		NotifiedAt: time.Now(),
	}

	err = repositories.CreateRemindNotiec(notice)
	if err != nil {
		logger.PrintErr("create remind notice", err)
		return err
	}

	return nil
}
