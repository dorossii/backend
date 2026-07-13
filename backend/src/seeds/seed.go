package seeds

import (
	"backend/logger"
	"backend/models"

	"gorm.io/gorm/clause"
)


func Seed() error {
	err := baseTaskSeed()
	if err != nil {
		logger.Println("Failed to seed tasks: ", err)
		return err
	}
	return nil
}


func baseTaskSeed() error {
	baseTasks := []models.BaseTask{
		{
			BaseID: "158a17f6-949b-4275-8706-3c2142e76a43",
			TaskName: "洗濯物を干す",
			Description: "洗濯機から取り出した衣類をハンガーにかけ、日当たりの良い場所に干す。",
			DifficultyLevel: 2, // 難易度は1から5の範囲で設定
			DueTime: 3,			// 日数単位の期限
			ImageFlag: true,
			Tags: models.TaskTagLaundry,
		},
		{
			BaseID: "e77bf6ac-d74e-426b-a241-fc90ec299c54",
			TaskName: "部屋の掃除機がけ",
			Description: "リビングおよび各寝室の床を掃除機で丁寧に吸引する。",
			DifficultyLevel: 3,
			DueTime: 5,
			ImageFlag: true,
			Tags: models.TaskTagCleaning,
		},
		{
			BaseID: "c1d221c3-7fd8-4f50-8da0-7490cabeada2",
			TaskName: "食器洗い",
			Description: "夕食後に溜まった皿、コップ、調理器具を洗って水切りカゴに置く。",
			DifficultyLevel: 2,
			DueTime: 5,
			ImageFlag: true,
			Tags: models.TaskTagCooking,
		},
		{
			BaseID: "0cba1ba6-82aa-42f4-ad08-f474b35f09d0",
			TaskName: "ゴミ出し",
			Description: "家中のゴミ箱を集め、指定の分別ルールに従ってゴミステーションへ運ぶ。",
			DifficultyLevel: 1,
			DueTime: 7,
			ImageFlag: true,
			Tags: models.TaskTagCleaning,
		},
		{
			BaseID: "a0d891d7-56f3-4685-b807-051ff89972b2",
			TaskName: "お風呂掃除",
			Description: "浴槽、床、蛇口周りをスポンジと洗剤を使って磨く。",
			DifficultyLevel: 5,
			DueTime: 14,
			ImageFlag: true,
			Tags: models.TaskTagCleaning,
		},
		{
			BaseID: "a6a51aee-ec3c-49dc-8191-a9544d31c6cf",
			TaskName: "トイレ掃除",
			Description: "便器内、便座、床を専用クリーナーで除菌・清掃する。",
			DifficultyLevel: 5,
			DueTime: 14,
			ImageFlag: false,
			Tags: models.TaskTagCleaning,
		},
		{
			BaseID: "798f8933-f759-4176-adb9-99371rc014087",
			TaskName: "窓拭き",
			Description: "リビングの大きな窓の汚れを拭き取り、視界をクリアにする。",
			DifficultyLevel: 4,
			DueTime: 14,
			ImageFlag: true,
			Tags: models.TaskTagCleaning,
		},
		{
			BaseID: "77dfe5bb-4818-4f20-84bc-70e71642d163",
			TaskName: "食材の買い出し",
			Description: "スーパーで翌日以降に必要な食材リストを確認しながら購入する。",
			DifficultyLevel: 3,
			DueTime: 3,
			ImageFlag: true,
			Tags: models.TaskTagAther,
		},
        {
            BaseID: "9c3e986b-12d4-48f8-b391-7d1a2e99f4a1",
            TaskName: "郵便物の整理",
            Description: "溜まっている郵便物やチラシを確認し、必要なものと不要なもの（ゴミ）に分ける。",
            DifficultyLevel: 1,
            DueTime: 7,
            ImageFlag: false,
            Tags: models.TaskTagAther,
        },
        {
            BaseID: "4b2f15a9-3d02-4c28-98e1-5e8f3b202c89",
            TaskName: "ベッドメイキング",
            Description: "布団を整え、枕やシーツの乱れを直す。",
            DifficultyLevel: 1,
            DueTime: 1,
            ImageFlag: true,
            Tags: models.TaskTagCleaning,
        },
        {
            BaseID: "8a7c5d9e-1f2b-4e63-b892-0c5d4f1a9e72",
            TaskName: "アイロンがけ",
            Description: "洗濯して乾いたシャツやハンカチのシワを伸ばし、綺麗に折りたたむ。",
            DifficultyLevel: 4,
            DueTime: 7,
            ImageFlag: true,
            Tags: models.TaskTagLaundry,
        },
        {
            BaseID: "d2e4f6a8-c0b2-4d1e-9f3a-5b7c9d1e3f5a",
            TaskName: "冷蔵庫の整理と賞味期限確認",
            Description: "冷蔵庫の中身を確認し、期限切れの食材がないかチェックして掃除する。",
            DifficultyLevel: 3,
            DueTime: 14,
            ImageFlag: true,
            Tags: models.TaskTagCooking,
        },
        {
            BaseID: "e5a3c9b7-1d8f-4e2a-8c6b-3f0d5a9b2e4c",
            TaskName: "玄関の掃除",
            Description: "玄関のたたきを掃き掃除し、靴を揃えて綺麗にする。",
            DifficultyLevel: 2,
            DueTime: 10,
            ImageFlag: true,
            Tags: models.TaskTagCleaning,
        },
	}	

	return models.DB.Create(&baseTasks).Clauses(clause.OnConflict{DoNothing: true}).Error

}

