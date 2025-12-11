package main

import (
	"fmt"
	"log"

	config "github.com/EmersonRabelo/first-api-go/internal/config"
	interfaces "github.com/EmersonRabelo/first-api-go/internal/config/interfaces"
	"github.com/EmersonRabelo/first-api-go/internal/controller"
	database "github.com/EmersonRabelo/first-api-go/internal/database"
	"github.com/EmersonRabelo/first-api-go/internal/repository"
	"github.com/EmersonRabelo/first-api-go/internal/router"
	"github.com/EmersonRabelo/first-api-go/internal/service"
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

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := controller.NewUserHandler(userService)

	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository)
	postHandler := controller.NewPostHandler(postService)

	r := router.SetupRouter(userHandler, postHandler)

	port := setting.GetServer().Port

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Falha ao iniciar servidor:", err)
	}
}
