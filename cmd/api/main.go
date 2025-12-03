package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	config "github.com/EmersonRabelo/first-api-go/internal/config"
	interfaces "github.com/EmersonRabelo/first-api-go/internal/config/interfaces"
	database "github.com/EmersonRabelo/first-api-go/internal/database"
	"github.com/gin-gonic/gin"
)

var setting interfaces.SettingProvider

func init() {
	fmt.Println("Application initializing...")

	setting = config.GetSetting()

	config.InitDatabase()

	fmt.Println("Initialized.")
}

func main() {

	db := config.GetDB()
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Falha ao executar migrations:", err)
	}

	fmt.Println(setting)
}

func healthCheck(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"status": "running",
		"time":   time.Now(),
	})
}
