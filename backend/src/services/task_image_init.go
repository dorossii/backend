package services

import (
	"errors"
	"os"
)

var uploadDir string

func InitTaskImageService() error {
	uploadDir = os.Getenv("TASK_IMAGE_DIR")

	if uploadDir == "" {
		return errors.New("TASK_IMAGE_DIR is not set")
	}

	return os.MkdirAll(uploadDir, 0755)
}