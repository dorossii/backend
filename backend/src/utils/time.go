package utils

import (
	"log"
	"time"
)

// JST はアプリケーション全体で使う日本標準時のタイムゾーン
var jst *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Printf("failed to load Asia/Tokyo location, falling back to fixed offset: %v", err)
		loc = time.FixedZone("Asia/Tokyo", 9*60*60)
	}
	jst = loc
}

// GetJST は日本標準時(Asia/Tokyo)の*time.Locationを返す
func GetJST() *time.Location {
	return jst
}

// NowJST は現在時刻を日本標準時(JST)で返す
func NowJST() time.Time {
	return time.Now().In(jst)
}
