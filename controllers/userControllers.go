package controllers

import (
	"mygram-api/database"
	"mygram-api/helpers"
	"mygram-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

// RegisterUser godoc
// @Summary register new user
// @Description register new user
// @Tags users
// @Accept json
// @Produce json
// @Param models.User body models.User true "create user"
// @Success 200 {object} models.User
// @Router /users/register [post]
func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	User := models.User{}
	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	err := db.Debug().Create(&User).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       User.ID,
		"username": User.Username,
		"email":    User.Email,
		"Age":      User.Age,
	})
}

// LoginUser godoc
// @Summary login user
// @Description login user
// @Tags users
// @Accept json
// @Produce json
// @Param models.User body models.User true "create user"
// @Success 200 {object} models.User
// @Router /users/login [post]
func UserLogin(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType

	User := models.User{}
	password := ""

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	password = User.Password
	err := db.Debug().Where("username=? OR email=?", User.Username, User.Email).Take(&User).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email/password",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))
	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email password",
		})
		return
	}

	token := helpers.GenerateToken(User.ID, User.Email, User.Username, User.Age)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
