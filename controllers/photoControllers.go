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
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	photos := models.Photo{}
	userId := uint(userData["id"].(float64))
	err := db.Find(&photos, "user_id=?", userId).Error

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
			"error_message": fmt.Sprintf("photo with id %v not found", photoId),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photo": photo,
	})
}

func UploadPhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = contentType, db

	photo := models.Photo{}
	userId := uint(userData["id"].(float64))

	_, _ = photo, userId

	// process filePhoto
	photoFile, err := c.FormFile("photo")
	title := c.PostForm("title")
	caption := c.PostForm("caption")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Bad Request",
			"error_message": "No file is received",
		})
		return
	}

	// generate and save photoFile
	photoUrl, err := helpers.SaveFile(photoFile, c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_status":  "Internal Server Error",
			"error_message": "Unable to save the file",
		})
		return
	}

	// redefine (karena form-data)
	photo.UserID = userId
	photo.PhotoUrl = photoUrl
	photo.Title = title
	photo.Caption = caption
	err = db.Debug().Create(&photo).Error

	// generate data user
	user := models.User{}
	user.ID = uint(userData["id"].(float64))
	user.Username = string(userData["username"].(string))
	user.Email = string(userData["email"].(string))
	user.Age = uint(userData["age"].(float64))

	photo.User = &user

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_status":  "Bad Request",
			"error_message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, photo)
}

// func UpdatePhoto(c *gin.Context) {

// }

// func DeletePhoto(c *gin.Context) {

// }
