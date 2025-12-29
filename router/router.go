package router

import (
	"net/http"
	"time"

	"github.com/EmersonRabelo/first-api-go/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *controller.UserHandler, postHandler *controller.PostHandler, likeHandler *controller.LikeHandler, replyHandler *controller.ReplyHandler, reportHandler *controller.ReportHandler) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"status": "running",
				"time":   time.Now(),
			})
		})

		users := v1.Group("/users")
		{
			users.GET("", userHandler.FindAll)
			users.GET("/:id", userHandler.FindById)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}

		posts := v1.Group("/posts")
		{
			posts.GET("", postHandler.FindAll)
			posts.GET("/:id", postHandler.FindById)
			posts.POST("", postHandler.Create)
			posts.PUT("/:id", postHandler.Update)
			posts.DELETE("/:id", postHandler.Delete)

			posts.POST("/:id/report", reportHandler.Create)
		}

		like := v1.Group("/likes")
		{
			like.GET("", likeHandler.FindAll)
			like.GET("/:id", likeHandler.FindById)
			like.POST("", likeHandler.Create)
			like.DELETE("/:id", likeHandler.Delete)
		}

		reply := v1.Group("/replies")
		{
			reply.GET("", replyHandler.FindAll)
			reply.GET("/:id", replyHandler.FindById)
			reply.POST("", replyHandler.Create)
			reply.PUT("/:id", replyHandler.Update)
			reply.DELETE("/:id", replyHandler.Delete)
		}

		report := v1.Group("/reports")
		{
			report.GET("")
			report.GET("/:id")
		}

	}

	return r
}
