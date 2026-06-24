package services

import (
	"backend/logger"
	"backend/models"
	"backend/repositories"
	"net/http"

	"errors"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrTaskNotFound             = errors.New("タスクが見つかりません")
	ErrInvalidTaskStatus        = errors.New("無効なタスクステータスです")
	ErrInvalidRequest           = errors.New("必要なパラメータの不足です")
	ErrTaskExpired              = errors.New("タスクの有効期間外です")
	ErrTaskStatusAlreadyUpdated = errors.New("すでにタスクステータスが更新されています")
	ErrTaskPermissionDenied     = errors.New("タスクを操作する権限がありません")

	ErrUnsupportedImageType = errors.New("対応していない画像形式です")
	ErrEmptyImageFile       = errors.New("空の画像ファイルです")
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

// タスク写真アップロード
func PostUploadImage(userID string, taskID string, fileHeader *multipart.FileHeader) error {
	task, err := repositories.GetTask(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}

		return err
	}

	// 自分のタスクか確認
	if task.UserID != userID {
		return ErrTaskPermissionDenied
	}

	oldImageID := task.ImageID

	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 画像か判定
	header := make([]byte, 512)
	n, err := src.Read(header)
	if err != nil && err != io.EOF {
		return err
	}

	contentType := http.DetectContentType(header[:n])

	// 画像形式判定
	var ext string
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	default:
		return ErrUnsupportedImageType
	}

	// 読み取り位置を先頭へ戻す
	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	fileName := uuid.NewString() + ext

	dstPath := filepath.Join(
		uploadDir,
		fileName,
	)

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	written, err := io.Copy(dst, src)
	if err != nil {
		dst.Close()
		os.Remove(dstPath)
		return err
	}

	if err := dst.Close(); err != nil {
		os.Remove(dstPath)
		return err
	}

	if written == 0 {
		os.Remove(dstPath)
		return ErrEmptyImageFile
	}

	// DB更新
	err = repositories.UpdateTaskImage(taskID, fileName)
	if err != nil {
		os.Remove(dstPath)
		return err
	}

	// 古い画像削除
	if oldImageID != "" {

		oldPath := filepath.Join(
			uploadDir,
			oldImageID,
		)

		_ = os.Remove(oldPath)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return PutTaskStatusResponse{}, ErrTaskNotFound
		}
		return PutTaskStatusResponse{}, err
	}

	// タスク所有者 or フレンドのみ操作可能
	if task.UserID != userID {
		friendShip, err := repositories.GetFriendShipAny(
			userID,
			task.UserID,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return PutTaskStatusResponse{},
					ErrTaskPermissionDenied
			}

			return PutTaskStatusResponse{}, err
		}

		if friendShip.UserID == "" {
			return PutTaskStatusResponse{},
				ErrTaskPermissionDenied
		}
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
		return PutTaskStatusResponse{},
			ErrTaskStatusAlreadyUpdated
	}

	if task.Status == models.TaskStatusCompleted &&
		newStatus != models.TaskStatusCompleted {
		return PutTaskStatusResponse{},
			ErrTaskStatusAlreadyUpdated
	}

	// レスキューに設定されているユーザーの保存
	var rescueUserIDs []models.HelpTargets

	switch status {
	case TaskStatusComplete:
		tx := models.DB.Begin()
		defer tx.Rollback()

		// 完了処理
		baseTask, err := repositories.GetBaseTask(task.BaseID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return PutTaskStatusResponse{}, ErrTaskNotFound
			}
			return PutTaskStatusResponse{}, err
		}

		// TODO:渡してる数字がintなのは許して後修正する
		err = repositories.UpdateTaskStatus(tx, taskID, models.TaskStatusCompleted)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		user, err := repositories.GetUser(userID)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		difficultyLevel := baseTask.DifficultyLevel * GarbagePower // 汚さ数値の計算

		// 自分の汚さの更新
		err = repositories.UpdateDirtLevel(tx, userID, -difficultyLevel)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		// 嫌がらせ相手の選出
		targetUserID := user.TargetUser

		// レスキュー対象除外
		rescueUserIDs, err = repositories.GetRescueUserIDs(userID)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		rescueMap := make(map[string]bool)

		for _, id := range rescueUserIDs {
			rescueMap[id.FriendID] = true
		}

		if targetUserID != "" {
			// 指定ターゲットがレスキュー対象なら再抽選
			if rescueMap[targetUserID] {
				targetUserID = ""
			}
		}

		if targetUserID == "" {
			friends, err := repositories.GetFriends(userID)
			if err != nil {
				return PutTaskStatusResponse{}, err
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
			err = repositories.UpdateDirtLevel(tx, targetUserID, difficultyLevel)
			if err != nil {
				return PutTaskStatusResponse{}, err
			}

			// 汚した相手に通知処理
			notice := &models.TrashNotice{
				NoticeID:   uuid.NewString(),
				SenderID:   userID,
				ReceiverID: targetUserID,
				Count:      baseTask.DifficultyLevel, // 難易度=ゴミの数
				CreatedAt:  time.Time{},
			}

			err = repositories.CreateTrashNotice(tx, notice)
			if err != nil {
				return PutTaskStatusResponse{}, err
			}
		}

		// TODO: レスキュー処理
		var validRescueUsers []models.HelpTargets

		// 空レコードをユーザとして判定しないように
		for _, rescueUser := range rescueUserIDs {
			if rescueUser.FriendID != "" {
				validRescueUsers = append(validRescueUsers, rescueUser)
			}
		}

		if len(validRescueUsers) > 0 {
			reduceAmount := difficultyLevel / len(validRescueUsers)

			for _, rescueUser := range validRescueUsers {
				err := repositories.UpdateDirtLevel(
					tx,
					rescueUser.FriendID,
					-reduceAmount,
				)
				if err != nil {
					return PutTaskStatusResponse{}, err
				}
			}
		}

		if err := tx.Commit().Error; err != nil {
			logger.PrintErr("commit transaction", err)
			return PutTaskStatusResponse{}, err
		}

		return PutTaskStatusResponse{
			IsChanged:    true,
			RequireImage: false,
		}, nil

	case TaskStatusPending:
		// 認証待ち処理
		if task.RequireImage && task.ImageID == "" {
			return PutTaskStatusResponse{}, ErrInvalidRequest
		}

		err = repositories.UpdateTaskStatus(models.DB, taskID, models.TaskStatusPending)
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
		tx := models.DB.Begin()
		defer tx.Rollback()

		// 未完了処理
		if task.Status != models.TaskStatusPending {
			return PutTaskStatusResponse{}, ErrTaskStatusAlreadyUpdated
		}

		if message == "" {
			return PutTaskStatusResponse{}, ErrInvalidRequest
		}

		err = repositories.UpdateTaskStatus(tx, taskID, models.TaskStatusImcomplete)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		err = repositories.UpdateTaskMessage(tx, taskID, message)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		if err := tx.Commit().Error; err != nil {
			logger.PrintErr("commit transaction", err)
			return PutTaskStatusResponse{}, err
		}

		return PutTaskStatusResponse{
			IsChanged: true,
		}, nil
	}

	return PutTaskStatusResponse{}, nil
}

func GetTaskImage(imageID string) (filePath string, contentType string, err error) {
	candidates := []struct {
		ext         string
		contentType string
	}{
		{".jpg", "image/jpeg"},
		{".png", "image/png"},
	}

	for _, c := range candidates {
		path := filepath.Join(uploadDir, imageID+c.ext)
		if _, statErr := os.Stat(path); statErr == nil {
			return path, c.contentType, nil
		}
	}

	return "", "", ErrTaskNotFound
}
