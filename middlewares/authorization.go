package middlewares

import (
	"mygram-api/database"
	"mygram-api/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func PhotoAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		photoId, err := strconv.Atoi(c.Param("photoId"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_status":  "Bad Request",
				"error_message": "Invalid parameter",
			})
			return
		}

		userData := c.MustGet("userData").(jwt.MapClaims)
		userId := uint(userData["id"].(float64))
		photo := models.Photo{}

		err = db.Select("user_id").First(&photo, uint(photoId)).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error_status":  "Data not found",
				"error_message": "Photo doesn't exist",
			})
			return
		}

		if photo.UserID != userId {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_status":  "Unauthorized",
				"error_message": "Your are not allowed to access this data",
			})
			return
		}

		c.Next()
	}
}

func CommentAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		commentId, err := strconv.Atoi(c.Param("commentId"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_status":  "Bad Request",
				"error_message": "Invalid parameter",
			})
			return
		}

		userData := c.MustGet("userData").(jwt.MapClaims)
		userId := uint(userData["id"].(float64))
		comment := models.Comment{}

		err = db.Select("user_id").First(&comment, uint(commentId)).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error_status":  "Data not found",
				"error_message": "Comment doesn't exist",
			})
			return
		}

		if comment.UserID != userId {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_status":  "Unauthorized",
				"error_message": "You are not allowed to access this data",
			})
			return
		}

		c.Next()
	}
}

func SocmedAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		socmedId, err := strconv.Atoi(c.Param("socmedId"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_status":  "Bad Request",
				"error_message": "Invalid parameter",
			})
			return
		}

		userData := c.MustGet("userData").(jwt.MapClaims)
		userId := uint(userData["id"].(float64))
		socmed := models.SocialMedia{}

		err = db.Select("user_id").First(&socmed, uint(socmedId)).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error_status":  "Data not found",
				"error_message": "Socmed doesn't exist",
			})
			return
		}

		if socmed.UserID != userId {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_status":  "Unauthorized",
				"erorr_message": "You are not allowed to access thi data",
			})
			return
		}

		c.Next()
	}
}
