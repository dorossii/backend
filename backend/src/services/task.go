package services

import (
	"backend/logger"
	"backend/models"
	"backend/repositories"
	"backend/utils"
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
	ErrTaskNotFound         = errors.New("タスクが見つかりません")
	ErrInvalidTaskStatus    = errors.New("無効なタスクステータスです")
	ErrInvalidRequest       = errors.New("必要なパラメータの不足です")
	ErrTaskExpired          = errors.New("タスクの有効期間外です")
	ErrInvalidTaskState     = errors.New("不正な状態遷移です")
	ErrTaskPermissionDenied = errors.New("タスクを操作する権限がありません")

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

// TaskPendingResponse は GET /user/tasks/pending のレスポンス要素
type TaskPendingResponse struct {
	TaskID      string `json:"taskId"`
	UserID      string `json:"userId"`
	TaskName    string `json:"taskName"`
	Tag         int    `json:"tag"`
	Description string `json:"description"`
	StartDate   int64  `json:"startDate"`
	EndTime     int64  `json:"endTime"`
	ImageID     string `json:"imageId"`
}

// GetPendingTasks はログインユーザーのフレンドが行った承認待ちタスク一覧を取得する
func GetPendingTasks(userID string) ([]TaskPendingResponse, error) {
	friends, err := repositories.GetFriends(userID)
	if err != nil {
		logger.PrintErr("フレンド一覧の取得に失敗", "userID", userID, "error", err)
		return []TaskPendingResponse{}, err
	}

	friendIDs := make([]string, 0, len(friends))
	for _, f := range friends {
		friendIDs = append(friendIDs, f.UserID)
	}

	rows, err := repositories.GetPendingTasksForFriends(friendIDs)
	if err != nil {
		logger.PrintErr("承認待ちタスクの取得に失敗", "userID", userID, "error", err)
		return []TaskPendingResponse{}, err
	}

	result := make([]TaskPendingResponse, 0, len(rows))
	for _, row := range rows {
		result = append(result, TaskPendingResponse{
			TaskID:      row.TaskID,
			UserID:      row.UserID,
			TaskName:    row.TaskName,
			Tag:         row.Tag,
			Description: row.Description,
			StartDate:   row.StartDate.Unix(),
			EndTime:     row.EndTime.Unix(),
			ImageID:     row.ImageID,
		})
	}

	return result, nil
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
		NotifiedAt: utils.NowJST(),
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
		return models.TaskStatusIncomplete, nil

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

	// タスクを更新できるのか検証
	task, err := validateTaskStatusUpdate(userID, taskID, status)
	if err != nil {
		return PutTaskStatusResponse{}, err
	}

	newStatus, err := ParseTaskStatus(status)
	if err != nil {
		return PutTaskStatusResponse{}, err
	}

	// タスクのステータスを変更できる状況なのか検証
	actor := "owner"

	if task.UserID != userID {
		actor = "friend"
	}

	err = validateTaskTransition(task, actor, newStatus, message)
	if err != nil {
		return PutTaskStatusResponse{}, err
	}

	switch newStatus {
	case models.TaskStatusCompleted:
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

		err = repositories.UpdateTaskStatus(tx, taskID, models.TaskStatusCompleted)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}
		// レスキューに設定されているユーザーの保存
		var rescueUserIDs []models.HelpTargets

		rescueUserIDs, err = repositories.GetRescueUserIDs(userID)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		// 汚さ更新
		user, dirtAmount, err := applyTaskCompletionEffect(tx, userID, baseTask, rescueUserIDs)
		if err != nil {
			return PutTaskStatusResponse{}, err
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
			reduceAmount := dirtAmount / len(validRescueUsers)

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

		lastCompleted := user.LastTaskCompletedAt
		now := utils.NowJST()

		diffDays := -1

		if lastCompleted != nil {
			diffDays = int(
				now.Truncate(24*time.Hour).
					Sub(lastCompleted.Truncate(24*time.Hour)).
					Hours() / 24,
			)
		}

		newCombo := 1

		if diffDays == 0 { // すでに今日分のコンボは加算されている
			newCombo = user.Combo
		} else if diffDays == 1 { // 最終更新が昨日なのでコンボ追加
			newCombo = user.Combo + 1
		}

		err = repositories.UpdateUserCombo(tx, userID, newCombo, now)
		if err != nil {
			return PutTaskStatusResponse{}, err
		}

		if err := tx.Commit().Error; err != nil {
			logger.PrintErr("commit transaction", err)
			return PutTaskStatusResponse{}, err
		}

		return PutTaskStatusResponse{
			IsChanged:    true,
			RequireImage: false,
		}, nil

	case models.TaskStatusPending:
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

	case models.TaskStatusIncomplete:
		tx := models.DB.Begin()
		defer tx.Rollback()

		// 未完了処理
		err = repositories.UpdateTaskStatus(tx, taskID, models.TaskStatusIncomplete)
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

// 更新して良いか検証
func validateTaskStatusUpdate(userID, taskID, status string) (models.Task, error) {
	task, err := repositories.GetTask(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}

	// タスク所有者 or フレンドのみ操作可能
	if task.UserID != userID {
		friendShip, err := repositories.GetFriendShipAny(userID, task.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.Task{},
					ErrTaskPermissionDenied
			}

			return models.Task{}, err
		}

		if friendShip == nil || friendShip.UserID == "" {
			return models.Task{},
				ErrTaskPermissionDenied
		}
	}

	// タスクの有効期間 検証
	now := time.Now()
	if now.Before(task.StartTime) || now.After(task.EndTime) {
		return models.Task{}, ErrTaskExpired
	}

	return task, nil
}

type TransitionRule struct {
	Actor          string
	CurrentStatus  models.TaskStatus
	NextStatus     models.TaskStatus
	RequireImage   bool
	RequireMessage bool
}

// 遷移ルール
var taskTransitionRules = []TransitionRule{
	{
		Actor:         "owner",
		CurrentStatus: models.TaskStatusIncomplete,
		NextStatus:    models.TaskStatusCompleted,
		RequireImage:  false,
	},
	{
		Actor:         "owner",
		CurrentStatus: models.TaskStatusIncomplete,
		NextStatus:    models.TaskStatusPending,
		RequireImage:  true,
	},
	{
		Actor:         "friend",
		CurrentStatus: models.TaskStatusPending,
		NextStatus:    models.TaskStatusCompleted,
	},
	{
		Actor:          "friend",
		CurrentStatus:  models.TaskStatusPending,
		NextStatus:     models.TaskStatusIncomplete,
		RequireMessage: true,
	},
}

// 遷移ルール検証
func validateTaskTransition(task models.Task, actor string, nextStatus models.TaskStatus, message string) error {

	for _, rule := range taskTransitionRules {

		if rule.Actor != actor {
			continue
		}

		if rule.CurrentStatus != task.Status {
			continue
		}

		if rule.NextStatus != nextStatus {
			continue
		}

		// 画像必須
		if rule.RequireImage &&
			task.RequireImage &&
			task.ImageID == "" {
			return ErrInvalidRequest
		}

		// メッセージ必須
		if rule.RequireMessage &&
			message == "" {
			return ErrInvalidRequest
		}

		return nil
	}

	return ErrInvalidTaskState
}

// 汚さ更新
func applyTaskCompletionEffect(tx *gorm.DB, userID string, baseTask models.BaseTask, rescueUserIDs []models.HelpTargets) (models.User, int, error) {
	user, err := repositories.GetUser(userID)
	if err != nil {
		return models.User{}, 0, err
	}

	dirtAmount := baseTask.DifficultyLevel * GarbagePower

	// 自分の汚さ減少
	err = repositories.UpdateDirtLevel(tx, userID, -dirtAmount)
	if err != nil {
		return models.User{}, 0, err
	}

	// 嫌がらせ相手の選出
	targetUserID, err := chooseTrashTarget(userID, user.TargetUser, rescueUserIDs)
	if err != nil {
		return models.User{}, 0, err
	}

	if targetUserID == "" {
		return *user, dirtAmount, nil
	}

	// 相手の汚さ増加
	err = repositories.UpdateDirtLevel(
		tx,
		targetUserID,
		dirtAmount,
	)
	if err != nil {
		return models.User{}, 0, err
	}

	notice := &models.TrashNotice{
		NoticeID:   uuid.NewString(),
		SenderID:   userID,
		ReceiverID: targetUserID,
		Count:      baseTask.DifficultyLevel,
	}

	return *user, dirtAmount, repositories.CreateTrashNotice(tx, notice)
}

// 嫌がらせ相手
func chooseTrashTarget(userID string, targetUserID string, rescueUserIDs []models.HelpTargets) (string, error) {
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
			return "", err
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

	return targetUserID, nil
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

type PutMultiTasksStatusRequest struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func PutMultiTasksStatus(userID string, req []PutMultiTasksStatusRequest) error {
	tx := models.DB.Begin()
	defer tx.Rollback()

	var completedTasks []PutMultiTasksStatusRequest
	var pendingTasks []PutMultiTasksStatusRequest
	taskIDs := make(map[string]struct{})

	// ステータスごとに振り分け・タスクID重複除外
	for _, taskReq := range req {
		if _, exists := taskIDs[taskReq.ID]; exists {
			continue
		}

		taskIDs[taskReq.ID] = struct{}{}

		switch taskReq.Status {
		case TaskStatusComplete:
			completedTasks = append(completedTasks, taskReq)

		case TaskStatusPending:
			pendingTasks = append(pendingTasks, taskReq)

		default:
			return ErrInvalidTaskState
		}
	}

	// レスキュー対象ユーザーを取得
	var rescueUserIDs []models.HelpTargets

	if len(completedTasks) > 0 {
		var err error
		rescueUserIDs, err = repositories.GetRescueUserIDs(userID)
		if err != nil {
			return err
		}
	}

	// 完了タスクの検証・完了処理
	completedTaskIDs := make([]string, 0, len(completedTasks))
	totalDirt := 0

	for _, taskReq := range completedTasks {
		task, err := validateTaskStatusUpdate(userID, taskReq.ID, taskReq.Status)
		if err != nil {
			return err
		}

		// 複数更新はタスク所有者のみ
		if task.UserID != userID {
			return ErrTaskPermissionDenied
		}

		// 遷移可能か検証(""これなのはメッセージが必要ないから)
		if err := validateTaskTransition(task, "owner", models.TaskStatusCompleted, ""); err != nil {
			return err
		}

		baseTask, err := repositories.GetBaseTask(task.BaseID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTaskNotFound
			}
			return err
		}

		// 完了時の汚さ・嫌がらせ処理
		_, dirtAmount, err := applyTaskCompletionEffect(tx, userID, baseTask, rescueUserIDs)
		if err != nil {
			return err
		}

		completedTaskIDs = append(completedTaskIDs, taskReq.ID)
		totalDirt += dirtAmount
	}

	// 完了ステータスを一括更新
	if len(completedTaskIDs) > 0 {
		if err := repositories.UpdateTasksStatus(tx, completedTaskIDs, models.TaskStatusCompleted); err != nil {
			return err
		}
	}

	// Pendingタスクの検証
	pendingTaskIDs := make([]string, 0, len(pendingTasks))

	for _, taskReq := range pendingTasks {
		task, err := validateTaskStatusUpdate(userID, taskReq.ID, taskReq.Status)
		if err != nil {
			return err
		}

		// 複数更新はタスク所有者のみ
		if task.UserID != userID {
			return ErrTaskPermissionDenied
		}

		// 遷移可能か検証(""これなのはメッセージが必要ないから)
		if err := validateTaskTransition(task, "owner", models.TaskStatusPending, ""); err != nil {
			return err
		}

		pendingTaskIDs = append(pendingTaskIDs, taskReq.ID)
	}

	// Pendingステータスを一括更新
	if len(pendingTaskIDs) > 0 {
		if err := repositories.UpdateTasksStatus(tx, pendingTaskIDs, models.TaskStatusPending); err != nil {
			return err
		}
	}

	// 完了タスクがある場合の追加処理
	if len(completedTasks) > 0 {
		// レスキュー処理
		var validRescueUsers []models.HelpTargets

		for _, rescueUser := range rescueUserIDs {
			if rescueUser.FriendID != "" {
				validRescueUsers = append(validRescueUsers, rescueUser)
			}
		}

		if totalDirt > 0 && len(validRescueUsers) > 0 {
			reduceAmount := totalDirt / len(validRescueUsers)

			for _, rescueUser := range validRescueUsers {
				if err := repositories.UpdateDirtLevel(tx, rescueUser.FriendID, -reduceAmount); err != nil {
					return err
				}
			}
		}

		// コンボ更新
		user, err := repositories.GetUser(userID)
		if err != nil {
			return err
		}

		lastCompleted := user.LastTaskCompletedAt
		now := utils.NowJST()

		diffDays := -1

		if lastCompleted != nil {
			diffDays = int(now.Truncate(24*time.Hour).Sub(lastCompleted.Truncate(24*time.Hour)).Hours() / 24)
		}

		newCombo := 1

		if diffDays == 0 {
			newCombo = user.Combo
		} else if diffDays == 1 {
			newCombo = user.Combo + 1
		}

		if err := repositories.UpdateUserCombo(tx, userID, newCombo, now); err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.PrintErr("commit transaction", err)
		return err
	}

	return nil
}
