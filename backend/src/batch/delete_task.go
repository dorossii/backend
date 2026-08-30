package batch

import (
	"backend/models"
	"backend/utils"
	"time"
)

func DeleteTaskTicker() {
	go func() {
		// 次の午前0時(JST)までの待機時間を計算
		now := utils.NowJST()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, utils.GetJST())
		timer := time.NewTimer(next.Sub(now))
		defer timer.Stop()

		// 最初に次の午前0時まで待機
		<-timer.C
		DeleteTask()

		// 以降は24時間ごとに実行
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			DeleteTask()
		}
	}()
}

func DeleteTask() error {
	//期限切れのタスクを削除する
	if err := models.DB.Where("end_time < ?", utils.NowJST()).Delete(&models.Task{}).Error; err != nil {
		return err
	}
	return nil
}
