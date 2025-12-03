package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"status": "running",
				"time":   time.Now(),
			})
		})

	}

	return r
}
