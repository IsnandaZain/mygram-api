package router

import (
	"mygram-api/controllers"

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
	// photoRouter := r.Group("/photos")
	// {
	// 	photoRouter.GET("/", controllers.PhotoGetAll)
	// 	photoRouter.GET("/:photoId", controllers.PhotoGet)
	// 	photoRouter.POST("/", controllers.UploadPhoto)
	// 	photoRouter.PUT("/:photoId", controllers.UpdatePhoto)
	// 	photoRouter.DELETE("/:photoId", controllers.DeletePhoto)
	// }

	// // Comment
	// commentRouter := r.Group("/comments")
	// {
	// 	commentRouter.GET("/", controllers.CommentGetAll)
	// 	commentRouter.GET("/:commentId", controllers.CommentGet)
	// 	commentRouter.POST("/", controllers.PostComment)
	// 	commentRouter.PUT("/:commentId", controllers.UpdateComment)
	// 	commentRouter.DELETE("/:commentId", controllers.DeleteComment)
	// }

	// // Social Media
	// socmedRouter := r.Group("/socmed")
	// {
	// 	socmedRouter.GET("/", controllers.SocmedGetAll)
	// 	socmedRouter.GET("/:socmedId", controllers.SocmedGet)
	// 	socmedRouter.POST("/", controllers.AddSocmed)
	// 	socmedRouter.PUT("/:socmedId", controllers.UpdateSocmed)
	// 	socmedRouter.DELETE("/:socmedId", controllers.DeleteSocmed)
	// }

	return r
}
