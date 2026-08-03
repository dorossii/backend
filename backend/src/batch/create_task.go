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
		// 次の午前0時(JST)までの待機時間を計算
		now := utils.NowJST()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, utils.GetJST())
		timer := time.NewTimer(next.Sub(now))
		defer timer.Stop()

		// 最初に次の午前0時まで待機
		<-timer.C
		CreateTask()

		// 以降は24時間ごとに実行
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			CreateTask()
		}
	}()
}

func CreateTask() error {

	// 1ユーザーにつき何個タスクを作成するかの定数
	const tasksPerUser = 2

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
	if len(baseTasks) < tasksPerUser {
		return errors.New("insufficient base tasks available (minimum 2 required)")
	}

	// 過去に各ユーザーへ割り当て済みの BaseID を全件取得
	// （ステータスに関わらず、これまでに割り当てられたことがあるタスクを重複除外対象にする）
	var assignedRows []struct {
		UserID string
		BaseID string
	}
	if err := models.DB.Model(&models.Task{}).
		Select("user_id, base_id").
		Find(&assignedRows).Error; err != nil {
		return err
	}

	// ユーザーごとに「過去に割り当て済みの BaseID セット」を構築
	assignedByUser := make(map[string]map[string]bool, len(userIDs))
	for _, row := range assignedRows {
		if assignedByUser[row.UserID] == nil {
			assignedByUser[row.UserID] = make(map[string]bool)
		}
		assignedByUser[row.UserID][row.BaseID] = true
	}

	// 乱数生成器の初期化
	r := rand.New(rand.NewSource(utils.NowJST().UnixNano()))

	var tasksToInsert []models.Task
	now := utils.NowJST()

	// 各ユーザーに対してランダムに2つのタスクを選出してスライスに格納
	for _, userID := range userIDs {
		assignedSet := assignedByUser[userID]

		// まだ割り当てたことのない BaseTask のみを候補にする
		candidates := make([]models.BaseTask, 0, len(baseTasks))
		for _, bt := range baseTasks {
			if !assignedSet[bt.BaseID] {
				candidates = append(candidates, bt)
			}
		}

		// 未割り当ての候補が必要数に満たない場合、
		// 「全BaseTaskを消化した」とみなして候補を全BaseTaskにリセットする
		// （＝新しい周回として再び選出可能にする）
		if len(candidates) < tasksPerUser {
			candidates = baseTasks
		}

		// 同じユーザーに同じタスクが重複して割り当たらないようにインデックスをシャッフル
		shuffledIndices := r.Perm(len(candidates))

		// 上位2つのランダムなタスクを選択
		for i := 0; i < tasksPerUser; i++ {
			baseTask := candidates[shuffledIndices[i]]

			// DueTime（期限）の仕様に合わせて終了時間を計算
			endTime := now.Add(time.Duration(baseTask.DueTime) * 24 * time.Hour)

			// UUIDを生成
			uuid, err := utils.Genid()
			if err != nil {
				return err
			}
			// imageflagがtrueの要素の中で10%の確率でRequireImageをtrueにする
			requireImage := false
			if baseTask.ImageFlag && r.Float64() < 0.1 {
				requireImage = true
			}

			task := models.Task{
				TaskID:       uuid,
				BaseID:       baseTask.BaseID,
				UserID:       userID,
				Status:       models.TaskStatusIncomplete,
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
