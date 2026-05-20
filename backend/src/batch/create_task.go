package batch

import (
	"backend/models"
	"backend/utils"
	"errors"
	"math/rand"
	"time"
)

func CreateTaskTicker() {
	go func() {
		//24時間ごとにタスクを作成する
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				CreateTask()
			}
		}

	}()
}

func CreateTask() error {
	// 全ユーザーにランダムでタスクを作成する
	var userIDs []string
	if err := models.DB.Model(&models.User{}).Pluck("user_id", &userIDs).Error; err != nil {
		return err
	}

	if len(userIDs) == 0 {
		return nil // ユーザーがいない場合は何もしない
	}

	// 全てのBaseTaskを取得
	var baseTasks []models.BaseTask
	if err := models.DB.Find(&baseTasks).Error; err != nil {
		return err
	}

	// ベースタスクが2つ未満だと「1ユーザーにつき2タスク」を満たせないためエラーハンドリング
	if len(baseTasks) < 2 {
		return errors.New("insufficient base tasks available (minimum 2 required)")
	}

	// 乱数生成器の初期化
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var tasksToInsert []models.Task
	now := time.Now()

	// 各ユーザーに対してランダムに2つのタスクを選出してスライスに格納
	for _, userID := range userIDs {
		// 同じユーザーに同じタスクが重複して割り当たらないようにインデックスをシャッフル
		shuffledIndices := r.Perm(len(baseTasks))

		// 上位2つのランダムなタスクを選択
		for i := 0; i < 2; i++ {
			baseTask := baseTasks[shuffledIndices[i]]

			// DueTime（期限）の仕様に合わせて終了時間を計算 
			endTime := now.Add(time.Duration(baseTask.DueTime) * 24 * time.Hour)

			uuid, err := utils.Genid() // 重複を避けるためUUID等のユニークIDを生成
			if err != nil {
				return err
			}
			// imageflagがtrueの要素の中で10%の確率でRequireImageをtrueにする
			requireImage := false
			if baseTask.ImageFlag && r.Float64() < 0.1 {
				requireImage = true
			}

			task := models.Task{
				TaskID:       uuid, // 重複を避けるためUUID等のユニークIDを生成
				UserID:       userID,
				Status:       models.TaskStatusPending,
				StartTime:    now,
				EndTime:      endTime,
				ImageID:      "", // 初期状態は空
				RequireImage: requireImage,
			}
			tasksToInsert = append(tasksToInsert, task)
		}
	}

	// トランザクション内でバルクインサート（一括保存）を実行
	if err := models.DB.CreateInBatches(&tasksToInsert, 100).Error; err != nil {
		return err
	}

	return nil
}
