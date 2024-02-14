package controllers

import (
	"fmt"
	"mygram-api/database"
	"mygram-api/helpers"
	"mygram-api/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func PhotoGetAll(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	photos := models.Photo{}
	err := db.Find(&photos).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_status":  "Failed get data",
			"error_message": "Something wrong when try to get data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photos": photos,
	})
}

func PhotoGet(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	photo := models.Photo{}
	photoId, _ := strconv.Atoi(c.Param("photoId"))
	userId := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&photo)
	} else {
		c.ShouldBind(&photo)
	}

	photo.UserID = userId
	photo.ID = uint(photoId)

	err := db.First(&photo, "id = ?", photoId).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "Photo not found",
			"error_message": fmt.Sprintf("car with id %v not found", photoId),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photo": photo,
	})
}

// func UploadPhoto(c *gin.Context) {

// }

// func UpdatePhoto(c *gin.Context) {

// }

// func DeletePhoto(c *gin.Context) {

// }
