package router

import (
	"mygram-api/controllers"
	"mygram-api/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	// User
	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
	}

	// Photo
	photoRouter := r.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.GET("/", controllers.PhotoGetAll)
		photoRouter.GET("/:photoId", middlewares.PhotoAuthorization(), controllers.PhotoGet)
		photoRouter.POST("/", controllers.UploadPhoto)
		// photoRouter.PUT("/:photoId", controllers.UpdatePhoto)
		// photoRouter.DELETE("/:photoId", controllers.DeletePhoto)
	}

	// Comment
	commentRouter := r.Group("/comments")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.GET("/", controllers.CommentGetAll)
		commentRouter.GET("/:photoId", controllers.CommentGet)
		commentRouter.POST("/:photoId", controllers.PostComment)
		commentRouter.PUT("/:commentId", middlewares.CommentAuthorization(), controllers.UpdateComment)
		commentRouter.DELETE("/:commentId", middlewares.CommentAuthorization(), controllers.DeleteComment)
	}

	// // Social Media
	socmedRouter := r.Group("/socmed")
	{
		socmedRouter.GET("/", controllers.SocmedGetAll)
		socmedRouter.GET("/:socmedId", controllers.SocmedGet)
		socmedRouter.POST("/", controllers.PostSocmed)
		socmedRouter.PUT("/:socmedId", middlewares.SocmedAuthorization(), controllers.UpdateSocmed)
		socmedRouter.DELETE("/:socmedId", middlewares.SocmedAuthorization(), controllers.DeleteSocmed)
	}

	return r
}
