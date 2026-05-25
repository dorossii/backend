package batch

// バッチ処理のコントローラー
func Run() {
	
	//タスク作成
	CreateTaskTicker()

	//汚さ増加
	DirtTicker()

	//HP変化
	ChangeHPTicker()

	//タスク削除
	DeleteTaskTicker()
}
