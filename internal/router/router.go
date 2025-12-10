package router

import (
	"net/http"
	"time"

	"github.com/EmersonRabelo/first-api-go/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *controller.UserHandler) *gin.Engine {
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

	}

	return r
}
