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

func SocmedGetAll(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	socmeds := models.SocialMedia{}
	userId := uint(userData["id"].(float64))
	err := db.Find(&socmeds, "user_id=?", userId).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_status":  "Failed get data",
			"error_message": "Something wrong when try to get data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"socmeds": socmeds,
	})
}

func SocmedGet(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	socmed := models.SocialMedia{}
	socmedId, _ := strconv.Atoi(c.Param("socmedId"))

	err := db.First(&socmed, "socmed_id=?", socmedId).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "Socmed not found",
			"error_message": fmt.Sprintf("socmed with id %v not found", socmedId),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"socmed": socmed,
	})
}

func PostSocmed(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	newSocmed := models.SocialMedia{}
	userId := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&newSocmed)
	} else {
		c.ShouldBind(&newSocmed)
	}

	newSocmed.UserID = userId
	err := db.Debug().Create(&newSocmed).Error

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
	newSocmed.User = &user

	c.JSON(http.StatusCreated, newSocmed)
}

func UpdateSocmed(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	socmedUpdate := models.SocialMedia{}
	socmedId, _ := strconv.Atoi(c.Param("socmedId"))
	if err := c.ShouldBindJSON(&socmedUpdate); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Bad Request",
			"error_message": "Bad Request",
		})
		return
	}

	var socmed models.SocialMedia
	err := db.First(&socmed, "id=?", uint(socmedId)).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "Social Media not found",
			"error_message": fmt.Sprintf("social media with id %v not found", socmedId),
		})
		return
	}

	// update socmed
	db.Model(&socmed).Where("id=?", uint(socmedId)).Updates(models.SocialMedia{
		Name:           socmedUpdate.Name,
		SocialMediaUrl: socmedUpdate.SocialMediaUrl,
	})

	socmedUpdate.ID = socmed.ID
	socmedUpdate.UserID = socmed.UserID

	// generate data user
	user := models.User{}
	user.ID = uint(userData["id"].(float64))
	user.Username = string(userData["username"].(string))
	user.Email = string(userData["email"].(string))
	user.Age = uint(userData["age"].(float64))
	socmedUpdate.User = &user

	c.JSON(http.StatusCreated, gin.H{
		"socmed": socmedUpdate,
	})
}

func DeleteSocmed(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	socmeId, _ := strconv.Atoi(c.Param("socmedId"))
	socmed := models.SocialMedia{}
	err := db.First(&socmed, "id=?", uint(socmeId)).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_status":  "Data not found",
			"error_message": fmt.Sprintf("social media with id %v not found", socmeId),
		})
		return
	}

	db.Delete(&socmed)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("social media with id %v has been deleted", socmeId),
	})
}
