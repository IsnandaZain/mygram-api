package controllers

import (
	"fmt"
	"mygram-api/database"
	"mygram-api/helpers"
	"mygram-api/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// GetAllPhoto godoc
// @Summary Get All Photo from user
// @Description Get All Photo from user
// @Tags photos
// @Accept json
// @Produce json
// @Success 200 {object} models.Photo
// @Router /photos [get]
func PhotoGetAll(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	photos := []models.Photo{}
	userId := uint(userData["id"].(float64))
	err := db.Preload("Comments").Preload("User").Find(&photos, "user_id=?", userId).Error

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

// GetPhoto godoc
// @Summary Get photo details
// @Description Get details of one photo
// @Tags photos
// @Accept json
// @Produce json
// @Param Id path int true "ID for the photo"
// @Success 200 {object} models.Photo
// @Router /photos/{photoId} [get]
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

	err := db.Preload("Comments").Preload("User").First(&photo, "id=?", photoId).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "Photo not found",
			"error_message": fmt.Sprintf("photo with id %v not found", photoId),
		})
		return
	}

	photo.UserID = userId
	photo.ID = uint(photoId)

	c.JSON(http.StatusOK, gin.H{
		"photo": photo,
	})
}

// PostPhoto godoc
// @Summary Post photo
// @Description Post photo by input user
// @Tags photos
// @Accept json
// @Produce json
// @Param models.Photo body models.Photo true "create photo"
// @Success 200 {object} models.Photo
// @Router /photos [post]
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

// UpdatePhoto godoc
// @Summary Update photo for a given Id
// @Description Update detail of photo corresponding to the input Id
// @Tags photos
// @Accept json
// @Produce json
// @Param models.Photo body models.Photo true "update photo"
// @Success 200 {object} models.Photo
// @Router /photos [put]
func UpdatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	photoUpdate := models.Photo{}
	userId := uint(userData["id"].(float64))
	photoId, _ := strconv.Atoi(c.Param("photoId"))
	title := c.PostForm("title")
	caption := c.PostForm("caption")

	// update
	photoUpdate.Title = title
	photoUpdate.Caption = caption

	// make sure photo exist
	var photo models.Photo
	err := db.First(&photo, "id=?", photoId).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "Data not found",
			"error_message": fmt.Sprintf("photo with id %v not found", photoId),
		})
		return
	}

	// process file
	photoFile, errFile := c.FormFile("photo")
	var photoUrl string
	if errFile == nil { // terdapat foto
		// generate and save photoFile
		photoUrl, err = helpers.SaveFile(photoFile, c)

		// delete file lama
		_ = os.Remove("upload/" + strings.Split(photo.PhotoUrl, "/")[2])
	} else {
		photoUrl = photo.PhotoUrl
	}

	// update photo
	db.Model(&photo).Where("id=?", photoId).Updates(models.Photo{
		Title:    photoUpdate.Title,
		Caption:  photoUpdate.Caption,
		PhotoUrl: photoUrl,
	})

	// redefine
	photoUpdate.UserID = userId
	photoUpdate.PhotoUrl = photoUrl

	// generate data user
	user := models.User{}
	user.ID = uint(userData["id"].(float64))
	user.Username = string(userData["username"].(string))
	user.Email = string(userData["email"].(string))
	user.Age = uint(userData["age"].(float64))
	photoUpdate.User = &user

	c.JSON(http.StatusCreated, gin.H{
		"photo": photoUpdate,
	})
}

// DeletePhoto godoc
// @Summary Delete photo for a given Id
// @Description Delete photo corresponding to the param Id
// @Tags photos
// @Accept json
// @Produce json
// @Param Id path int tru "ID for the photo"
// @Success 200 "{'message': 'Photo has been deleted'}"
// @Router /photos/{photoId} [delete]
func DeletePhoto(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	photo := models.Photo{}

	err := db.First(&photo, "id=?", uint(photoId)).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_status":  "Data not found",
			"error_message": fmt.Sprintf("photo with id %v not found", photoId),
		})
		return
	}

	// delete comment first
	var comments models.Comment
	err = db.Find(&comments, "photo_id=?", photoId).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_status":  fmt.Sprintf("failed to process data with id %v", photoId),
			"error_message": fmt.Sprintf("failed to process data with id %v", photoId),
		})
		return
	}

	// delete
	db.Delete(&comments)
	db.Delete(&photo)

	// delete file
	_ = os.Remove("upload/" + strings.Split(photo.PhotoUrl, "/")[2])

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("photo with id %v has been deleted", photoId),
	})
}
