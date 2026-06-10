package services

import (
	"backend/logger"
	"backend/models"
	"backend/repositories"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTaskStatus        = errors.New("無効なタスクステータスです")
	ErrInvalidRequest           = errors.New("必要なパラメータの不足です")
	ErrTaskExpired              = errors.New("タスクの有効期間外です")
	ErrTaskNotFound             = errors.New("タスクが見つかりません")
	ErrTaskStatusAlreadyUpdated = errors.New("すでにタスクステータスが更新されています")
)

func GetTasks(userID string) ([]repositories.TaskResponse, error) {
	tasks, err := repositories.GetUserTasks(userID)
	if err != nil {
		logger.PrintErr("タスクの取得に失敗", "userID", userID, "error", err)
		return []repositories.TaskResponse{}, err
	}
	return tasks, nil
}

// 煽りメッセージの登録
func PostTaskTauntMessage(userId string, friendId string, msg string) error {
	// フレンド存在確認
	friendShip, err := repositories.GetFriendShipAny(userId, friendId)
	if err != nil {
		return err
	}

	if friendShip == nil {
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

type PutTaskStatusResponse struct {
	IsChanged    bool
	RequireImage bool
}

const (
	TaskStatusComplete   = "complete"
	TaskStatusPending    = "pending"
	TaskStatusIncomplete = "incomplete"
)

const GarbagePower = 18 // TODO: 難易度1 = 18p汚さ減る

// 文字列を TaskStatus に変換する関数
func ParseTaskStatus(s string) (models.TaskStatus, error) {
	switch s {
	case TaskStatusIncomplete:
		return models.TaskStatusImcomplete, nil

	case TaskStatusPending:
		return models.TaskStatusPending, nil

	case TaskStatusComplete:
		return models.TaskStatusCompleted, nil

	default:
		return 0, ErrInvalidTaskStatus
	}
}

// 　タスクステータス更新(完了•未完了)
func PutTaskStatus(userID, taskID, status, message string) (PutTaskStatusResponse, error) {
	task, err := repositories.GetTask(taskID)
	if err != nil {
		return PutTaskStatusResponse{}, err
	}

	if task.TaskID == "" {
		return PutTaskStatusResponse{}, ErrTaskNotFound
	}

	// タスクの有効期間 検証
	now := time.Now()
	if now.Before(task.StartTime) || now.After(task.EndTime) {
		return PutTaskStatusResponse{}, ErrTaskExpired
	}

	newStatus, err := ParseTaskStatus(status)
	if err != nil {
		return PutTaskStatusResponse{}, err
	}

	if task.Status == newStatus {
		return PutTaskStatusResponse{}, ErrTaskStatusAlreadyUpdated
	}

	switch status {
	case TaskStatusComplete:
		// 完了処理
		baseTask, err := repositories.GetBaseTask(task.BaseID)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		if baseTask.BaseID == "" {
			return PutTaskStatusResponse{}, ErrTaskNotFound
		}

		// TODO:渡してる数字がintなのは許して後修正する
		err = repositories.UpdateTaskStatus(taskID, 2)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		user, err := repositories.GetUser(userID)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		difficultyLevel := baseTask.DifficultyLevel * GarbagePower // 汚さ数値の計算

		// 自分の汚さの更新
		err = repositories.UpdateDirtLevel(userID, -difficultyLevel)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		// 嫌がらせ相手の選出
		targetUserID := user.TargetUser

		if targetUserID == "" {
			friends, err := repositories.GetFriends(userID)
			if err != nil {
				return PutTaskStatusResponse{}, err
			}

			// レスキュー対象除外
			rescueUserIDs, err := repositories.GetRescueUserIDs(userID)
			if err != nil {
				return PutTaskStatusResponse{}, err
			}

			rescueMap := make(map[string]bool)

			// 検索しやすい形(Map)に変換
			for _, id := range rescueUserIDs {
				rescueMap[id.FriendID] = true
			}

			var candidates []string

			for _, friend := range friends {
				if !rescueMap[friend.UserID] {
					candidates = append(candidates, friend.UserID)
				}
			}

			if len(candidates) > 0 {
				targetUserID = candidates[rand.Intn(len(candidates))]
			}
		}

		// 相手の汚さの更新
		if targetUserID != "" {
			err = repositories.UpdateDirtLevel(targetUserID, difficultyLevel)
			if err != nil {
				return PutTaskStatusResponse{}, err
			}
		}

		// TODO: 汚した相手に通知処理

		// TODO: レスキュー処理

		return PutTaskStatusResponse{
			IsChanged:    true,
			RequireImage: false,
		}, nil

	case TaskStatusPending:
		// 認証待ち処理
		if task.RequireImage && task.ImageID == "" {
			return PutTaskStatusResponse{}, ErrInvalidRequest
		}

		err = repositories.UpdateTaskStatus(taskID, 1)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		if task.RequireImage {
			return PutTaskStatusResponse{
				IsChanged:    true,
				RequireImage: true,
			}, nil
		}

		return PutTaskStatusResponse{
			IsChanged:    true,
			RequireImage: false,
		}, nil

	case TaskStatusIncomplete:
		// 未完了処理
		if task.Status != models.TaskStatusPending {
			return PutTaskStatusResponse{}, ErrTaskStatusAlreadyUpdated
		}

		if message == "" {
			return PutTaskStatusResponse{}, ErrInvalidRequest
		}

		err = repositories.UpdateTaskStatus(taskID, 0)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		err = repositories.UpdateTaskMessage(taskID, message)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		return PutTaskStatusResponse{
			IsChanged: true,
		}, nil

	default:
		return PutTaskStatusResponse{},
			ErrInvalidTaskStatus
	}
}
