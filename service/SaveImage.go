package service

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"mime/multipart"
	"os"
)

func SaveImage(file multipart.File, fileType string) (string, string) {
	dirPath := "./blog_image/"
	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.Mkdir(dirPath, 0755)
		}
	}

	fn := uuid.NewV4().String() + fileType
	newFile, err := os.Create(dirPath + fn)
	if err != nil {
		log.Println(err)
		return "", "上传图片失败"
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		log.Println(err)
		return "", "上传图片失败"
	}

	return fn, ""
}
