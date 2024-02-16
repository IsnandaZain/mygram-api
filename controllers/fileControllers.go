package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetFileUpload(c *gin.Context) {
	filename := c.Param("filename")
	fileBytes, err := ioutil.ReadFile("upload/" + filename)

	ext := strings.Split(filename, ".")[1]
	var datatype string
	if ext == "jpg" || ext == "png" {
		datatype = "image/" + ext
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_status":  "Failed to get data",
			"error_message": "Something wrong when try to get data",
		})
		return
	}

	c.Data(http.StatusOK, datatype, fileBytes)
}
