package oss

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestUpload(t *testing.T) {
	godotenv.Load("../../.env")
	imgurl := "https://api.telegram.org/file/bot6712916653:AAGjQhZIXJ6ZBO901P0tSX91mL1-liPM4TQ/photos/file_74.jpg"
	uploader := NewImageUploader()
	_, err := uploader.Upload(imgurl)
	if err != nil {
		t.Fatal(err)
	}
}
