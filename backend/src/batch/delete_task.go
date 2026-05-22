package batch

import (
	"backend/models"
	"time"
)

func DeleteTaskTicker() {
	go func() {
		//24時間ごとにタスクを削除する
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				DeleteTask()
			}	
		}
	}()
}

func DeleteTask() error {
	//期限切れのタスクを削除する
	if err := models.DB.Where("end_time < ?", time.Now()).Delete(&models.Task{}).Error; err != nil {
		return err
	}
	return nil
}
