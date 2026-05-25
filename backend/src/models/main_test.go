package models_test

import (
	"backend/models"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	models.InitForTest()
	os.Exit(m.Run())
}
