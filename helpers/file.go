package helpers

import (
	"mime/multipart"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	BASEURL = "127.0.0.1:8080"
)

func SaveFile(fileData *multipart.FileHeader, c *gin.Context) (filename string, err error) {

	// generate filename
	ext := filepath.Ext(fileData.Filename)
	newPhotoName := uuid.New().String() + ext

	if err := c.SaveUploadedFile(fileData, "upload/"+newPhotoName); err != nil {
		return "", err
	}

	// generate url file
	photoUrl := GenerateUrlFile(newPhotoName)
	return photoUrl, nil
}

func GenerateUrlFile(filename string) (url string) {
	url = BASEURL + "/fileUpload/" + filename
	return
}
