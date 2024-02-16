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

// GetAllComment godoc
// @Summary Get All Comment from user
// @Description Get All Comment from user
// @Tags comments
// @Accept json
// @Produce json
// @Success 200 {object} models.Comment
// @Router /comments [get]
func CommentGetAll(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	comments := []models.Comment{}
	userId := uint(userData["id"].(float64))
	err := db.Preload("User").Preload("Photo").Find(&comments, "user_id=?", userId).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_status":  "Failed get data",
			"error_message": "Something wrong when try to get data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})

}

// GetComment godoc
// @Summary Get all comment for photo
// @Description Get all comment for photo
// @Tags comments
// @Accept json
// @Produce json
// @Param Id path int true "ID for the photo"
// @Success 200 {object} models.Comment
// @Router /comments/{photoId} [get]
func CommentGet(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	comment := []models.Comment{}
	photoId, _ := strconv.Atoi(c.Param("photoId"))

	err := db.Preload("User").Preload("Photo").First(&comment, "photo_id = ?", photoId).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "Comment not found",
			"error_message": fmt.Sprintf("comment for photo with id %v not found", photoId),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}

// PostComment godoc
// @Summary Post comment for specific photo
// @Description Post comment for specific photo
// @Tags comments
// @Accept json
// @Produce json
// @Param models.Comment body models.Comment true "create comment"
// @Success 200 {object} models.Comment
// @Router /comments [post]
func PostComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	newComment := models.Comment{}
	userId := uint(userData["id"].(float64))
	photoId, _ := strconv.Atoi(c.Param("photoId"))
	if contentType == appJSON {
		c.ShouldBindJSON(&newComment)
	} else {
		c.ShouldBind(&newComment)
	}

	newComment.UserID = userId
	newComment.PhotoID = uint(photoId)
	err := db.Debug().Create(&newComment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_status":  "Bad Request",
			"error_message": "Bad Request",
		})
		return
	}

	// generate data user
	user := models.User{}
	user.ID = uint(userData["id"].(float64))
	user.Username = string(userData["username"].(string))
	user.Email = string(userData["email"].(string))
	user.Age = uint(userData["age"].(float64))
	newComment.User = &user

	// generate data photo
	var photoComment models.Photo
	err = db.First(&photoComment, "id=?", uint(photoId)).Error

	newComment.Photo = &photoComment

	c.JSON(http.StatusCreated, newComment)
}

// UpdateComment godoc
// @Summary Update comment for a given Id
// @Description Update comment corresponding to the input Id
// @Tags comments
// @Accept json
// @Produce json
// @Param models.Comment body models.Comment true "update comemnt"
// @Success 200 {object} models.Comment
// @Router /comments/{commentId} [put]
func UpdateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	commentUpdate := models.Comment{}
	commentId, _ := strconv.Atoi(c.Param("commentId"))
	if err := c.ShouldBindJSON(&commentUpdate); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Bad Request",
			"error_message": "Bad Request",
		})
		return
	}

	var comment models.Comment
	err := db.First(&comment, "id=?", uint(commentId)).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "Comment not found",
			"error_message": fmt.Sprintf("comment with id %v not found", commentId),
		})
		return
	}

	// update comment
	db.Model(&comment).Where("id=?", uint(commentId)).Updates(models.Comment{
		Message: commentUpdate.Message,
	})

	commentUpdate.ID = comment.ID
	commentUpdate.UserID = comment.UserID
	commentUpdate.PhotoID = comment.PhotoID

	// generate data user
	user := models.User{}
	user.ID = uint(userData["id"].(float64))
	user.Username = string(userData["username"].(string))
	user.Email = string(userData["email"].(string))
	user.Age = uint(userData["age"].(float64))
	commentUpdate.User = &user

	// generate data photo
	var photoComment models.Photo
	err = db.First(&photoComment, "id=?", comment.PhotoID).Error

	commentUpdate.Photo = &photoComment

	c.JSON(http.StatusCreated, gin.H{
		"comment": commentUpdate,
	})
}

// DeleteComment godoc
// @Summary Delete comment for a given Id
// @Description Delete comment corresponding to the param Id
// @Tags comments
// @Accept json
// @Produce json
// @Param Id path int tru "ID for the comment"
// @Success 200 "{'message': 'Comment has been deleted'}"
// @Router /comments/{commentId} [delete]
func DeleteComment(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	comment := models.Comment{}
	err := db.First(&comment, "id=?", uint(commentId)).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_status":  "Data not found",
			"error_message": fmt.Sprintf("comment with id %v not found", commentId),
		})
		return
	}

	db.Delete(&comment)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("comment with id %v has been deleted", commentId),
	})
}
