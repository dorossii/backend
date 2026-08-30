package services_test

import (
	"backend/batch"
	"backend/models"
	"backend/services"
	"backend/utils"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
	"errors"
)

// services.GetTasksのテスト
func TestGetTasks(t *testing.T) {
	// テスト用のユーザー作成
	if err := CreateSampleUser(); err != nil {
		t.Fatalf("テストユーザーの作成に失敗: %v", err)
	}

	// ベースタスクの準備(DueTimeは日数単位)
	baseTasks := []models.BaseTask{
		{
			BaseID:          "base-001",
			TaskName:        "部屋掃除",
			DueTime:         1,
			ImageId:         "image-001",
			ImageFlag:       true,
			Description:     "掃除して部屋をきれいにしよう",
			DifficultyLevel: 2,
			Tags:            0,
		},
		{
			BaseID:          "base-002",
			TaskName:        "洗濯物を干す",
			DueTime:         2,
			ImageId:         "image-002",
			ImageFlag:       false,
			Description:     "洗濯物を干すのを忘れないようにしよう",
			DifficultyLevel: 4,
			Tags:            3,
		},
		{
			BaseID:          "base-003",
			TaskName:        "夕飯を作る",
			DueTime:         3,
			ImageId:         "image-003",
			ImageFlag:       false,
			Description:     "夕飯を作ることを忘れないようにしよう",
			DifficultyLevel: 1,
			Tags:            2,
		},
	}
	if err := models.DB.Create(&baseTasks).Error; err != nil {
		t.Fatalf("failed to create dummy base tasks: %v", err)
	}

	err := batch.CreateTask()
	if err != nil {
		t.Fatalf("タスクの作成に失敗: %v", err)
	}

	// タスクを取得
	tasks, err := services.GetTasks("user-010")
	if err != nil {
		t.Fatalf("タスクの取得に失敗: %v", err)
	}

	// 取得したタスクの数を確認
	if len(tasks) == 0 {
		t.Errorf("タスクが見つかりません")
	}

	// 取得したタスクの難易度が1から5の範囲内であることを確認
	for _, task := range tasks {
		if task.DifficultyLevel < 1 || task.DifficultyLevel > 5 {
			t.Errorf("タスクの難易度が不正です: %d", task.DifficultyLevel)
		}
	}

	log.Printf("取得したタスク数: %d", len(tasks))

	// 取得したタスクの内容を確認
	for _, task := range tasks {
		log.Printf("タスクID: %s, タスク名: %s, 期限: %s,難易度: %d, 画像ID: %s", task.TaskID, task.TaskName, time.Unix(task.EndTime, 0).Format("2006-01-02"), task.DifficultyLevel, task.ImageID)
	}
}

// タスクテーブルを空にする（他のテストが作成したタスクの影響を排除する）
func truncateTasks(t *testing.T) {
	t.Helper()

	if err := models.DB.Exec("TRUNCATE TABLE tasks").Error; err != nil {
		t.Fatal(err)
	}
}

// services.GetPendingTasksのテスト
func TestGetPendingTasks(t *testing.T) {
	truncateFriendShips(t)
	truncateUsersAndRooms(t)
	truncateTasks(t)

	TestRegisterUser(t)
	// RegisterUser経由で自動生成されたタスクは本テストの検証対象外のため除去する
	truncateTasks(t)

	friend := models.FriendShips{
		UserID:   "user-001",
		FriendID: "user-002",
		Status:   models.FriendStatusAccepted,
	}
	if err := models.DB.Create(&friend).Error; err != nil {
		t.Fatal(err)
	}

	baseTask := models.BaseTask{
		BaseID:          "base-pending-001",
		TaskName:        "部屋掃除",
		DueTime:         1,
		ImageFlag:       true,
		Description:     "掃除して部屋をきれいにしよう",
		DifficultyLevel: 2,
		Tags:            0,
	}
	if err := models.DB.Create(&baseTask).Error; err != nil {
		t.Fatalf("failed to create dummy base task: %v", err)
	}

	pendingTask := models.Task{
		TaskID:    "task-pending-001",
		BaseID:    baseTask.BaseID,
		UserID:    "user-002",
		Status:    models.TaskStatusPending,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
		ImageID:   "img-001",
	}
	if err := models.DB.Create(&pendingTask).Error; err != nil {
		t.Fatalf("failed to create dummy pending task: %v", err)
	}

	// user-002自身が所有する未完了タスク（承認待ちではないので含まれないはず）
	incompleteTask := models.Task{
		TaskID:    "task-incomplete-001",
		BaseID:    baseTask.BaseID,
		UserID:    "user-002",
		Status:    models.TaskStatusIncomplete,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}
	if err := models.DB.Create(&incompleteTask).Error; err != nil {
		t.Fatalf("failed to create dummy incomplete task: %v", err)
	}

	tasks, err := services.GetPendingTasks("user-001")
	if err != nil {
		t.Fatalf("GetPendingTasks failed: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 pending task, got %d", len(tasks))
	}

	if tasks[0].TaskID != "task-pending-001" {
		t.Fatalf("unexpected taskId: %s", tasks[0].TaskID)
	}
	if tasks[0].UserID != "user-002" {
		t.Fatalf("unexpected userId: %s", tasks[0].UserID)
	}

	// フレンドがいないユーザーは空配列（nilではない）が返る
	emptyTasks, err := services.GetPendingTasks("user-003")
	if err != nil {
		t.Fatalf("GetPendingTasks failed: %v", err)
	}
	if emptyTasks == nil {
		t.Fatalf("expected empty slice, got nil")
	}
	if len(emptyTasks) != 0 {
		t.Fatalf("expected 0 pending tasks, got %d", len(emptyTasks))
	}
}

// フレンドがpending状態のタスクを持っていない場合（incomplete/completeのみ）は空配列が返る
func TestGetPendingTasks_FriendHasNoPendingTask(t *testing.T) {
	truncateFriendShips(t)
	truncateUsersAndRooms(t)
	truncateTasks(t)

	TestRegisterUser(t)
	// RegisterUser経由で自動生成されたタスクは本テストの検証対象外のため除去する
	truncateTasks(t)

	friend := models.FriendShips{
		UserID:   "user-001",
		FriendID: "user-002",
		Status:   models.FriendStatusAccepted,
	}
	if err := models.DB.Create(&friend).Error; err != nil {
		t.Fatal(err)
	}

	baseTask := models.BaseTask{
		BaseID:          "base-no-pending",
		TaskName:        "部屋掃除",
		DueTime:         1,
		ImageFlag:       true,
		Description:     "掃除して部屋をきれいにしよう",
		DifficultyLevel: 2,
		Tags:            0,
	}
	if err := models.DB.Create(&baseTask).Error; err != nil {
		t.Fatalf("failed to create dummy base task: %v", err)
	}

	// フレンドの未完了タスク（承認待ちではない）
	incompleteTask := models.Task{
		TaskID:    "task-incomplete-only",
		BaseID:    baseTask.BaseID,
		UserID:    "user-002",
		Status:    models.TaskStatusIncomplete,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}
	if err := models.DB.Create(&incompleteTask).Error; err != nil {
		t.Fatalf("failed to create dummy incomplete task: %v", err)
	}

	// フレンドの完了済みタスク
	completeTask := models.Task{
		TaskID:    "task-complete-only",
		BaseID:    baseTask.BaseID,
		UserID:    "user-002",
		Status:    models.TaskStatusCompleted,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}
	if err := models.DB.Create(&completeTask).Error; err != nil {
		t.Fatalf("failed to create dummy complete task: %v", err)
	}

	tasks, err := services.GetPendingTasks("user-001")
	if err != nil {
		t.Fatalf("GetPendingTasks failed: %v", err)
	}
	if tasks == nil {
		t.Fatalf("expected empty slice, got nil")
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 pending tasks, got %d", len(tasks))
	}
}

// フレンドではないユーザーの承認待ちタスクは取得できない
func TestGetPendingTasks_NonFriendTaskExcluded(t *testing.T) {
	truncateFriendShips(t)
	truncateUsersAndRooms(t)
	truncateTasks(t)

	TestRegisterUser(t)
	// RegisterUser経由で自動生成されたタスクは本テストの検証対象外のため除去する
	truncateTasks(t)

	// user-001とuser-002はフレンドだが、user-003とはフレンドではない
	friend := models.FriendShips{
		UserID:   "user-001",
		FriendID: "user-002",
		Status:   models.FriendStatusAccepted,
	}
	if err := models.DB.Create(&friend).Error; err != nil {
		t.Fatal(err)
	}

	baseTask := models.BaseTask{
		BaseID:          "base-non-friend",
		TaskName:        "部屋掃除",
		DueTime:         1,
		ImageFlag:       true,
		Description:     "掃除して部屋をきれいにしよう",
		DifficultyLevel: 2,
		Tags:            0,
	}
	if err := models.DB.Create(&baseTask).Error; err != nil {
		t.Fatalf("failed to create dummy base task: %v", err)
	}

	// フレンドではないuser-003の承認待ちタスク
	nonFriendPendingTask := models.Task{
		TaskID:    "task-non-friend-pending",
		BaseID:    baseTask.BaseID,
		UserID:    "user-003",
		Status:    models.TaskStatusPending,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
		ImageID:   "img-non-friend",
	}
	if err := models.DB.Create(&nonFriendPendingTask).Error; err != nil {
		t.Fatalf("failed to create dummy pending task: %v", err)
	}

	tasks, err := services.GetPendingTasks("user-001")
	if err != nil {
		t.Fatalf("GetPendingTasks failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 pending tasks (non-friend task must be excluded), got %d", len(tasks))
	}
}

// フレンド申請中（Pending）や拒否済み（Rejected）の relationshipは友達とみなされず、タスクも取得されない
func TestGetPendingTasks_FriendShipNotAcceptedExcluded(t *testing.T) {
	truncateFriendShips(t)
	truncateUsersAndRooms(t)
	truncateTasks(t)

	TestRegisterUser(t)

	// user-001 <-> user-002 は申請中
	pendingFriendShip := models.FriendShips{
		UserID:   "user-001",
		FriendID: "user-002",
		Status:   models.FriendStatusPending,
	}
	if err := models.DB.Create(&pendingFriendShip).Error; err != nil {
		t.Fatal(err)
	}

	// user-001 <-> user-003 は拒否済み
	rejectedFriendShip := models.FriendShips{
		UserID:   "user-001",
		FriendID: "user-003",
		Status:   models.FriendStatusRejected,
	}
	if err := models.DB.Create(&rejectedFriendShip).Error; err != nil {
		t.Fatal(err)
	}

	baseTask := models.BaseTask{
		BaseID:          "base-not-accepted",
		TaskName:        "部屋掃除",
		DueTime:         1,
		ImageFlag:       true,
		Description:     "掃除して部屋をきれいにしよう",
		DifficultyLevel: 2,
		Tags:            0,
	}
	if err := models.DB.Create(&baseTask).Error; err != nil {
		t.Fatalf("failed to create dummy base task: %v", err)
	}

	pendingUserTask := models.Task{
		TaskID:    "task-pending-friendship-pending",
		BaseID:    baseTask.BaseID,
		UserID:    "user-002",
		Status:    models.TaskStatusPending,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
		ImageID:   "img-002",
	}
	if err := models.DB.Create(&pendingUserTask).Error; err != nil {
		t.Fatalf("failed to create dummy pending task: %v", err)
	}

	rejectedUserTask := models.Task{
		TaskID:    "task-pending-friendship-rejected",
		BaseID:    baseTask.BaseID,
		UserID:    "user-003",
		Status:    models.TaskStatusPending,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
		ImageID:   "img-003",
	}
	if err := models.DB.Create(&rejectedUserTask).Error; err != nil {
		t.Fatalf("failed to create dummy pending task: %v", err)
	}

	tasks, err := services.GetPendingTasks("user-001")
	if err != nil {
		t.Fatalf("GetPendingTasks failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 pending tasks (pending/rejected friendships must be excluded), got %d", len(tasks))
	}
}

// フレンドシップが逆方向（friend_id = 自分）で登録されていても友達とみなされ、タスクが取得できる
func TestGetPendingTasks_ReverseFriendShipDirection(t *testing.T) {
	truncateFriendShips(t)
	truncateUsersAndRooms(t)
	truncateTasks(t)

	TestRegisterUser(t)
	// RegisterUser経由で自動生成されたタスクは本テストの検証対象外のため除去する
	truncateTasks(t)

	// user-002がuser-001に対して申請し承認された体（friend_id側が自分）
	friend := models.FriendShips{
		UserID:   "user-002",
		FriendID: "user-001",
		Status:   models.FriendStatusAccepted,
	}
	if err := models.DB.Create(&friend).Error; err != nil {
		t.Fatal(err)
	}

	baseTask := models.BaseTask{
		BaseID:          "base-reverse",
		TaskName:        "部屋掃除",
		DueTime:         1,
		ImageFlag:       true,
		Description:     "掃除して部屋をきれいにしよう",
		DifficultyLevel: 2,
		Tags:            0,
	}
	if err := models.DB.Create(&baseTask).Error; err != nil {
		t.Fatalf("failed to create dummy base task: %v", err)
	}

	pendingTask := models.Task{
		TaskID:    "task-reverse-pending",
		BaseID:    baseTask.BaseID,
		UserID:    "user-002",
		Status:    models.TaskStatusPending,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
		ImageID:   "img-reverse",
	}
	if err := models.DB.Create(&pendingTask).Error; err != nil {
		t.Fatalf("failed to create dummy pending task: %v", err)
	}

	tasks, err := services.GetPendingTasks("user-001")
	if err != nil {
		t.Fatalf("GetPendingTasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 pending task via reverse friendship, got %d", len(tasks))
	}
	if tasks[0].TaskID != "task-reverse-pending" {
		t.Fatalf("unexpected taskId: %s", tasks[0].TaskID)
	}
}

// 複数フレンドのうち、承認待ちタスクを持つフレンドの分だけ正しく返る
func TestGetPendingTasks_MultipleFriendsMixedStatus(t *testing.T) {
	truncateFriendShips(t)
	truncateUsersAndRooms(t)
	truncateTasks(t)

	TestRegisterUser(t)
	// RegisterUser経由で自動生成されたタスクは本テストの検証対象外のため除去する
	truncateTasks(t)

	// user-001はuser-002、user-003の両方とフレンド
	friendWith002 := models.FriendShips{
		UserID:   "user-001",
		FriendID: "user-002",
		Status:   models.FriendStatusAccepted,
	}
	if err := models.DB.Create(&friendWith002).Error; err != nil {
		t.Fatal(err)
	}
	friendWith003 := models.FriendShips{
		UserID:   "user-001",
		FriendID: "user-003",
		Status:   models.FriendStatusAccepted,
	}
	if err := models.DB.Create(&friendWith003).Error; err != nil {
		t.Fatal(err)
	}

	baseTask := models.BaseTask{
		BaseID:          "base-mixed",
		TaskName:        "部屋掃除",
		DueTime:         1,
		ImageFlag:       true,
		Description:     "掃除して部屋をきれいにしよう",
		DifficultyLevel: 2,
		Tags:            0,
	}
	if err := models.DB.Create(&baseTask).Error; err != nil {
		t.Fatalf("failed to create dummy base task: %v", err)
	}

	// user-002は承認待ちタスクを持つ
	pendingTask := models.Task{
		TaskID:    "task-mixed-pending",
		BaseID:    baseTask.BaseID,
		UserID:    "user-002",
		Status:    models.TaskStatusPending,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
		ImageID:   "img-mixed",
	}
	if err := models.DB.Create(&pendingTask).Error; err != nil {
		t.Fatalf("failed to create dummy pending task: %v", err)
	}

	// user-003は未完了タスクのみ（承認待ちではない）
	incompleteTask := models.Task{
		TaskID:    "task-mixed-incomplete",
		BaseID:    baseTask.BaseID,
		UserID:    "user-003",
		Status:    models.TaskStatusIncomplete,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}
	if err := models.DB.Create(&incompleteTask).Error; err != nil {
		t.Fatalf("failed to create dummy incomplete task: %v", err)
	}

	tasks, err := services.GetPendingTasks("user-001")
	if err != nil {
		t.Fatalf("GetPendingTasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 pending task, got %d", len(tasks))
	}
	if tasks[0].TaskID != "task-mixed-pending" {
		t.Fatalf("unexpected taskId: %s", tasks[0].TaskID)
	}
	if tasks[0].UserID != "user-002" {
		t.Fatalf("unexpected userId: %s", tasks[0].UserID)
	}
}

// レスポンスの各フィールドが正しくマッピングされていることを確認する
func TestGetPendingTasks_ResponseFieldsMapping(t *testing.T) {
	truncateFriendShips(t)
	truncateUsersAndRooms(t)
	truncateTasks(t)

	TestRegisterUser(t)
	// RegisterUser経由で自動生成されたタスクは本テストの検証対象外のため除去する
	truncateTasks(t)

	friend := models.FriendShips{
		UserID:   "user-001",
		FriendID: "user-002",
		Status:   models.FriendStatusAccepted,
	}
	if err := models.DB.Create(&friend).Error; err != nil {
		t.Fatal(err)
	}

	baseTask := models.BaseTask{
		BaseID:          "base-fields",
		TaskName:        "洗濯物を干す",
		DueTime:         2,
		ImageFlag:       false,
		Description:     "洗濯物を干すのを忘れないようにしよう",
		DifficultyLevel: 4,
		Tags:            3,
	}
	if err := models.DB.Create(&baseTask).Error; err != nil {
		t.Fatalf("failed to create dummy base task: %v", err)
	}

	startTime := time.Now().Truncate(time.Second)
	endTime := startTime.Add(24 * time.Hour)

	pendingTask := models.Task{
		TaskID:    "task-fields",
		BaseID:    baseTask.BaseID,
		UserID:    "user-002",
		Status:    models.TaskStatusPending,
		StartTime: startTime,
		EndTime:   endTime,
		ImageID:   "img-fields",
	}
	if err := models.DB.Create(&pendingTask).Error; err != nil {
		t.Fatalf("failed to create dummy pending task: %v", err)
	}

	tasks, err := services.GetPendingTasks("user-001")
	if err != nil {
		t.Fatalf("GetPendingTasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 pending task, got %d", len(tasks))
	}

	got := tasks[0]
	if got.TaskID != "task-fields" {
		t.Fatalf("unexpected TaskID: %s", got.TaskID)
	}
	if got.UserID != "user-002" {
		t.Fatalf("unexpected UserID: %s", got.UserID)
	}
	if got.TaskName != "洗濯物を干す" {
		t.Fatalf("unexpected TaskName: %s", got.TaskName)
	}
	if got.Tag != 3 {
		t.Fatalf("unexpected Tag: %d", got.Tag)
	}
	if got.Description != "洗濯物を干すのを忘れないようにしよう" {
		t.Fatalf("unexpected Description: %s", got.Description)
	}
	if got.StartDate != startTime.Unix() {
		t.Fatalf("unexpected StartDate: got %d, want %d", got.StartDate, startTime.Unix())
	}
	if got.EndTime != endTime.Unix() {
		t.Fatalf("unexpected EndTime: got %d, want %d", got.EndTime, endTime.Unix())
	}
	if got.ImageID != "img-fields" {
		t.Fatalf("unexpected ImageID: %s", got.ImageID)
	}
}

// テスト用のユーザーを作成する関数
func CreateSampleUser() error {
	users := []models.User{
		{
			UserID:      "user-010",
			UserName:    "syatyo",
			BirthDate:   time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
			Mailadress:  "user1@example.com",
			HealthPoint: 1000,
			DirtLevel:   0,
			Combo:       0,
			BgColor:     "#ffb6c1",
		},
	}
	if err := models.DB.Create(&users).Error; err != nil {
		return fmt.Errorf("failed to create dummy users: %v", err)
	}
	return nil
}

func truncateRemindNotices(t *testing.T) {
	t.Helper()

	if err := models.DB.Exec("TRUNCATE TABLE remind_notices").Error; err != nil {
		t.Fatal(err)
	}
}

func TestPostTaskTauntMessage(t *testing.T) {
	truncateFriendShips(t)
	truncateRemindNotices(t)

	TestRegisterUser(t)
	// seedFriend(t) TODO:呼び出し元がマージされてないので一時的に…
	friend := models.FriendShips{
		UserID:   "user-001",
		FriendID: "user-002",
		Status:   models.FriendStatusAccepted,
	}

	if err := models.DB.Create(&friend).Error; err != nil {
		t.Fatal(err)
	}

	err := services.PostTaskTauntMessage(
		"user-001",
		"user-002",
		"お前の部屋きたなすぎ",
	)

	if err != nil {
		t.Fatalf(
			"PostTaskTauntMessage failed: %v",
			err,
		)
	}

	var notice models.RemindNotice

	err = models.DB.
		First(&notice, "user_id = ?", "user-001").
		Error

	if err != nil {
		t.Fatalf(
			"record not found: %v",
			err,
		)
	}

	if notice.SenderID != "user-002" {
		t.Fatalf(
			"unexpected sender id: %s",
			notice.SenderID,
		)
	}

	if notice.Title != "お前の部屋きたなすぎ" {
		t.Fatalf(
			"unexpected title: %s",
			notice.Title,
		)
	}
}

func TestPostTaskTauntMessage_FriendNotFound(t *testing.T) {
	truncateFriendShips(t)
	truncateRemindNotices(t)

	TestRegisterUser(t)

	err := services.PostTaskTauntMessage(
		"user-001",
		"user-999",
		"test message",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrFriendNotFound {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

// 画像アップロードテスト用初期化
func setupUploadTest(t *testing.T) {
	t.Helper()

	t.Setenv(
		"TASK_IMAGE_DIR",
		"../assets/test-images",
	)

	if err := services.InitTaskImageService(); err != nil {
		t.Fatal(err)
	}
}

// 画像アップロードテストに利用するヘルパー
func createFileHeader(t *testing.T, path string) *multipart.FileHeader {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(
		"file",
		filepath.Base(path),
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(part, file); err != nil {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/upload",
		body,
	)

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	_, fileHeader, err := req.FormFile("file")
	if err != nil {
		t.Fatal(err)
	}

	return fileHeader
}

// 画像アップロード(正常系)
func TestPostUploadImage_JPEG(t *testing.T) {
	setupUploadTest(t)

	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-jpeg",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusIncomplete,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	fileHeader := createFileHeader(
		t,
		"../assets/test-images/test.jpg",
	)

	err := services.PostUploadImage(
		"user-001",
		"task-jpeg",
		fileHeader,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var updatedTask models.Task

	if err := models.DB.
		First(&updatedTask, "task_id = ?", "task-jpeg").
		Error; err != nil {
		t.Fatal(err)
	}

	if updatedTask.ImageID == "" {
		t.Fatal("expected image id")
	}
}

// 画像アップロード(画像差し替え)
func TestPostUploadImage_ReplaceImage(t *testing.T) {
	setupUploadTest(t)

	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-replace-image",
		BaseID:    "base-001",
		UserID:    "user-001",
		ImageID:   "old-image.jpg",
		Status:    models.TaskStatusIncomplete,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	fileHeader := createFileHeader(
		t,
		"../assets/test-images/test.jpg",
	)

	err := services.PostUploadImage(
		"user-001",
		"task-replace-image",
		fileHeader,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var updatedTask models.Task

	if err := models.DB.
		First(&updatedTask, "task_id = ?", "task-replace-image").
		Error; err != nil {
		t.Fatal(err)
	}

	if updatedTask.ImageID == "" {
		t.Fatal("expected image id")
	}

	if updatedTask.ImageID == "old-image.jpg" {
		t.Fatal("image was not replaced")
	}
}

// 画像アップロード(異常系：タスク不在)
func TestPostUploadImage_TaskNotFound(t *testing.T) {
	setupUploadTest(t)

	fileHeader := createFileHeader(
		t,
		"../assets/test-images/test.jpg",
	)

	err := services.PostUploadImage(
		"user-001",
		"not-found-task",
		fileHeader,
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrTaskNotFound {
		t.Fatalf("unexpected error: %v", err)
	}
}

// 画像アップロード(異常系：他人のタスク)
func TestPostUploadImage_PermissionDenied(t *testing.T) {
	setupUploadTest(t)

	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-other-user",
		BaseID:    "base-001",
		UserID:    "user-999",
		Status:    models.TaskStatusIncomplete,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	fileHeader := createFileHeader(
		t,
		"../assets/test-images/test.jpg",
	)

	err := services.PostUploadImage(
		"user-001",
		"task-other-user",
		fileHeader,
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrTaskPermissionDenied {
		t.Fatalf("unexpected error: %v", err)
	}
}

// 画像アップロード(異常系：JPEG/PNG以外)
func TestPostUploadImage_UnsupportedImageType(t *testing.T) {
	setupUploadTest(t)

	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-invalid-image",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusIncomplete,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	fileHeader := createFileHeader(
		t,
		"../assets/test-images/test.txt",
	)

	err := services.PostUploadImage(
		"user-001",
		"task-invalid-image",
		fileHeader,
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err != services.ErrUnsupportedImageType {
		t.Fatalf("unexpected error: %v", err)
	}
}

// タスクステータス更新(完了: 正常系)
func TestPutTaskStatus_Complete(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-Complete",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusIncomplete,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	resp, err := services.PutTaskStatus(
		"user-001",
		"task-Complete",
		services.TaskStatusComplete,
		"",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.IsChanged {
		t.Fatal("expected IsChanged=true")
	}

	var updatedTask models.Task

	if err := models.DB.
		First(&updatedTask, "task_id = ?", "task-Complete").
		Error; err != nil {
		t.Fatal(err)
	}

	if updatedTask.Status != models.TaskStatusCompleted {
		t.Fatalf(
			"unexpected status: %v",
			updatedTask.Status,
		)
	}
}

// タスクステータス更新(完了: タスク不存在)
func TestPutTaskStatus_TaskNotFound(t *testing.T) {
	TestRegisterUser(t)

	_, err := services.PutTaskStatus(
		"user-001",
		"task-TaskNotFound",
		services.TaskStatusComplete,
		"",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrTaskNotFound {
		t.Fatalf("unexpected error: %v", err)
	}
}

// タスクステータス更新(完了: 有効期間外)
func TestPutTaskStatus_Expired(t *testing.T) {
	CreateSampleUser()

	task := models.Task{
		TaskID:    "task-Expired",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusIncomplete,
		StartTime: utils.NowJST().Add(-2 * time.Hour),
		EndTime:   utils.NowJST().Add(-1 * time.Hour), // 既に終了
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-001",
		"task-Expired",
		services.TaskStatusComplete,
		"",
	)

	if err != services.ErrTaskExpired {
		t.Fatalf(
			"expected ErrTaskExpired, got %v",
			err,
		)
	}
}

// タスクステータス更新(認証待ち: 正常系)
func TestPutTaskStatus_Pending(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:       "task-pending",
		BaseID:       "base-001",
		UserID:       "user-001",
		Status:       models.TaskStatusIncomplete,
		StartTime:    utils.NowJST().Add(-1 * time.Hour),
		EndTime:      utils.NowJST().Add(1 * time.Hour),
		ImageID:      "image-001",
		RequireImage: false,
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	resp, err := services.PutTaskStatus(
		"user-001",
		"task-pending",
		services.TaskStatusPending,
		"",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.IsChanged {
		t.Fatal("expected IsChanged=true")
	}

	if resp.RequireImage {
		t.Fatal("expected RequireImage=false")
	}

	var updatedTask models.Task

	if err := models.DB.
		First(&updatedTask, "task_id = ?", "task-pending").
		Error; err != nil {
		t.Fatal(err)
	}

	if updatedTask.Status != models.TaskStatusPending {
		t.Fatalf(
			"unexpected status: %v",
			updatedTask.Status,
		)
	}
}

// タスクステータス更新(承認待ち: 写真必須なのに画像なし)
func TestPutTaskStatus_Pending_RequireImageButNoImageID(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:       "task-pending-no-image",
		BaseID:       "base-001",
		UserID:       "user-001",
		Status:       models.TaskStatusIncomplete,
		StartTime:    utils.NowJST().Add(-1 * time.Hour),
		EndTime:      utils.NowJST().Add(1 * time.Hour),
		ImageID:      "",   // 画像なし
		RequireImage: true, // 画像必須
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-001",
		"task-pending-no-image",
		services.TaskStatusPending,
		"",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrInvalidRequest {
		t.Fatalf("unexpected error: %v", err)
	}
}

// タスクステータス更新(未完了: 正常系)
func TestPutTaskStatus_Incomplete(t *testing.T) {
	TestRegisterUser(t)
	seedFriend(t)

	task := models.Task{
		TaskID:    "task-incomplete",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusPending, // Pendingから差し戻し
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	message := "普通に汚い"

	resp, err := services.PutTaskStatus(
		"user-002",
		"task-incomplete",
		services.TaskStatusIncomplete,
		message,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.IsChanged {
		t.Fatal("expected IsChanged=true")
	}

	var updatedTask models.Task

	if err := models.DB.
		First(&updatedTask, "task_id = ?", "task-incomplete").
		Error; err != nil {
		t.Fatal(err)
	}

	if updatedTask.Status != models.TaskStatusIncomplete {
		t.Fatalf(
			"unexpected status: %v",
			updatedTask.Status,
		)
	}

	if updatedTask.Message != message {
		t.Fatalf(
			"unexpected message: %v",
			updatedTask.Message,
		)
	}
}

// タスクステータス更新(未完了: 認証待ち以外から戻そうとする)
func TestPutTaskStatus_Incomplete_NotPending(t *testing.T) {
	TestRegisterUser(t)
	seedFriend(t)

	task := models.Task{
		TaskID:    "task-incomplete-not-pending",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusCompleted,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-002",
		"task-incomplete-not-pending",
		services.TaskStatusIncomplete,
		"差し戻し理由",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrInvalidTaskState {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

// タスクステータス更新(未完了: 拒否理由なし)
func TestPutTaskStatus_Incomplete_EmptyMessage(t *testing.T) {
	TestRegisterUser(t)
	seedFriend(t)

	task := models.Task{
		TaskID:    "task-incomplete-empty-message",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusPending,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-002",
		"task-incomplete-empty-message",
		services.TaskStatusIncomplete,
		"",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrInvalidRequest {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

// タスクステータス更新(ステータス不正)
func TestPutTaskStatus_InvalidStatus(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-invalid-status",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusIncomplete,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-001",
		"task-invalid-status",
		"invalid-status",
		"",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrInvalidTaskStatus {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

// タスクステータス更新(同じステータスへの更新完了)
func TestPutTaskStatus_AlreadyUpdated(t *testing.T) {
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-already-updated",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusPending,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	_, err := services.PutTaskStatus(
		"user-001",
		"task-already-updated",
		services.TaskStatusPending,
		"",
	)

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrInvalidTaskState {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

// 画像取得(正常系: JPEG)
func TestGetTaskImage_JPEG(t *testing.T) {
	setupUploadTest(t)
	TestRegisterUser(t)

	task := models.Task{
		TaskID:    "task-get-jpeg",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusPending,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	fileHeader := createFileHeader(t, "../assets/test-images/test.jpg")
	if err := services.PostUploadImage("user-001", "task-get-jpeg", fileHeader); err != nil {
		t.Fatalf("画像アップロードに失敗: %v", err)
	}

	var updatedTask models.Task
	if err := models.DB.First(&updatedTask, "task_id = ?", "task-get-jpeg").Error; err != nil {
		t.Fatal(err)
	}

	// ImageID から拡張子を除いた UUID 部分を取り出す
	imageUUID := updatedTask.ImageID[:len(updatedTask.ImageID)-len(filepath.Ext(updatedTask.ImageID))]

	filePath, contentType, err := services.GetTaskImage(imageUUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if filePath == "" {
		t.Fatal("expected non-empty filePath")
	}

	if contentType != "image/jpeg" {
		t.Fatalf("unexpected contentType: %s", contentType)
	}
}

// 画像取得(異常系: 存在しない imageId)
func TestGetTaskImage_NotFound(t *testing.T) {
	setupUploadTest(t)

	_, _, err := services.GetTaskImage("non-existent-uuid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err != services.ErrTaskNotFound {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPutMultiTasksStatus_Complete(t *testing.T) {
	TestRegisterUser(t)

	task1 := models.Task{
		TaskID:    "task-MultiComplete-001",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusIncomplete,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	task2 := models.Task{
		TaskID:    "task-MultiComplete-002",
		BaseID:    "base-001",
		UserID:    "user-001",
		Status:    models.TaskStatusIncomplete,
		StartTime: utils.NowJST().Add(-1 * time.Hour),
		EndTime:   utils.NowJST().Add(1 * time.Hour),
	}

	if err := models.DB.Create(&task1).Error; err != nil {
		t.Fatal(err)
	}

	if err := models.DB.Create(&task2).Error; err != nil {
		t.Fatal(err)
	}

	req := []services.PutMultiTasksStatusRequest{
		{
			ID:     task1.TaskID,
			Status: services.TaskStatusComplete,
		},
		{
			ID:     task2.TaskID,
			Status: services.TaskStatusPending,
		},
	}

	err := services.PutMultiTasksStatus("user-001", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var updatedTask1 models.Task
	if err := models.DB.First(&updatedTask1, "task_id = ?", task1.TaskID).Error; err != nil {
		t.Fatal(err)
	}

	var updatedTask2 models.Task
	if err := models.DB.First(&updatedTask2, "task_id = ?", task2.TaskID).Error; err != nil {
		t.Fatal(err)
	}

	if updatedTask1.Status != models.TaskStatusCompleted {
		t.Fatalf("unexpected status: %v", updatedTask1.Status)
	}

	if updatedTask2.Status != models.TaskStatusPending {
		t.Fatalf("unexpected status: %v", updatedTask2.Status)
	}
}

func TestPutMultiTasksStatus_Error(t *testing.T) {
	t.Run("存在しないタスク", func(t *testing.T) {
		TestRegisterUser(t)

		status := services.TaskStatusComplete

		req := []services.PutMultiTasksStatusRequest{
			{
				ID:     "not-found-mulit-task",
				Status: status,
			},
		}

		err := services.PutMultiTasksStatus("user-001", req)
		if !errors.Is(err, services.ErrTaskNotFound) {
			t.Fatalf("expected ErrTaskNotFound, got %v", err)
		}
	})

	t.Run("他ユーザーのタスク", func(t *testing.T) {
		TestRegisterUser(t)

		task := models.Task{
			TaskID:    "mulit-task-other-user",
			BaseID:    "base-001",
			UserID:    "user-002",
			Status:    models.TaskStatusIncomplete,
			StartTime: utils.NowJST().Add(-1 * time.Hour),
			EndTime:   utils.NowJST().Add(1 * time.Hour),
		}

		if err := models.DB.Create(&task).Error; err != nil {
			t.Fatal(err)
		}

		status := services.TaskStatusComplete

		req := []services.PutMultiTasksStatusRequest{
			{
				ID:     task.TaskID,
				Status: status,
			},
		}

		err := services.PutMultiTasksStatus("user-001", req)
		if !errors.Is(err, services.ErrTaskPermissionDenied) {
			t.Fatalf("expected ErrTaskPermissionDenied, got %v", err)
		}
	})

	t.Run("期限切れタスク", func(t *testing.T) {
		TestRegisterUser(t)

		task := models.Task{
			TaskID:    "mulit-task-expired",
			BaseID:    "base-001",
			UserID:    "user-001",
			Status:    models.TaskStatusIncomplete,
			StartTime: utils.NowJST().Add(-2 * time.Hour),
			EndTime:   utils.NowJST().Add(-1 * time.Hour),
		}

		if err := models.DB.Create(&task).Error; err != nil {
			t.Fatal(err)
		}

		status := services.TaskStatusComplete

		req := []services.PutMultiTasksStatusRequest{
			{
				ID:     task.TaskID,
				Status: status,
			},
		}

		err := services.PutMultiTasksStatus("user-001", req)
		if !errors.Is(err, services.ErrTaskExpired) {
			t.Fatalf("expected ErrTaskExpired, got %v", err)
		}
	})
}
