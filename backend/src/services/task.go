package services

import (
	"backend/logger"
	"backend/models"
	"backend/repositories"
	"net/http"

	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTaskNotFound             = errors.New("タスクが見つかりません")

	ErrTaskPermissionDenied = errors.New("このタスクを操作する権限がありません")
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
		return err
	}

	if task.TaskID == "" {
		return ErrTaskNotFound
	}

	// 自分のタスクか確認
	if task.UserID != userID {
		return ErrTaskPermissionDenied
	}

	oldImageID := task.ImageID

	uploadDir := "../assets/task-images"

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return err
	}

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
	defer dst.Close()

	written, err := io.Copy(dst, src)
	if err != nil {
		os.Remove(dstPath)
		return err
	}

	if written == 0 {
		os.Remove(dstPath)
		return errors.New("空ファイルです")
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
