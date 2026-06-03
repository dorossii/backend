package services

import (
	"backend/logger"
	"backend/models"
	"backend/repositories"
	"errors"
	"math/rand"
	"time"
)

var (
	ErrInvalidTaskStatus        = errors.New("無効なタスクステータスです")
	ErrImageRequired            = errors.New("画像必須です")
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

const GarbagePower = 18 // 難易度1 = 18p汚さ減る

// 　タスクステータス更新(完了•未完了)
func PutTaskStatus(userID, taskID, status, message string) (PutTaskStatusResponse, error) {
	task, err := repositories.GetTask(taskID)
	if err != nil {
		return PutTaskStatusResponse{}, err
	}

	if status == string(task.Status) {
		return PutTaskStatusResponse{}, ErrTaskStatusAlreadyUpdated
	}

	switch status {
	case TaskStatusComplete:
		// 完了処理
		baseTask, err := repositories.GetBaseTask(task.BaseID)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		// TODO: baseと分けるべき？？
		if task.TaskID == "" || baseTask.BaseID == "" {
			return PutTaskStatusResponse{}, ErrTaskNotFound
		}

		// タスクの有効期間 検証
		now := time.Now()
		if now.Before(task.StartTime) || now.After(task.EndTime) {
			return PutTaskStatusResponse{}, ErrTaskExpired
		}

		err = repositories.UpdateTaskStatus(taskID, TaskStatusComplete)
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
				if !rescueMap[friend.FriendID] {
					candidates = append(candidates, friend.FriendID)
				}
			}

			if len(friends) > 0 {
				targetUserID = candidates[rand.Intn(len(friends))]
			}
		}

		// 相手の汚さの更新
		if targetUserID != "" {
			err = repositories.UpdateDirtLevel(targetUserID, difficultyLevel)
			if err != nil {
				return PutTaskStatusResponse{}, err
			}
		}

		return PutTaskStatusResponse{
			IsChanged:    true,
			RequireImage: false,
		}, nil

	case TaskStatusPending:
		// 認証待ち処理
		if task.RequireImage && task.ImageID == "" {
			return PutTaskStatusResponse{}, ErrImageRequired
		}

		err = repositories.UpdateTaskStatus(taskID, TaskStatusPending)
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

		return PutTaskStatusResponse{
			IsChanged: true,
		}, nil

	default:
		return PutTaskStatusResponse{},
			ErrInvalidTaskStatus
	}
}
