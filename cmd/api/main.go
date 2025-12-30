package main

import (
	"fmt"
	"log"

	config "github.com/EmersonRabelo/first-api-go/internal/config"
	controller "github.com/EmersonRabelo/first-api-go/internal/controller"
	database "github.com/EmersonRabelo/first-api-go/internal/database"
	"github.com/EmersonRabelo/first-api-go/internal/queue"
	redis "github.com/EmersonRabelo/first-api-go/internal/redis"
	repository "github.com/EmersonRabelo/first-api-go/internal/repository"
	service "github.com/EmersonRabelo/first-api-go/internal/service"
	reportService "github.com/EmersonRabelo/first-api-go/internal/service/report"
	router "github.com/EmersonRabelo/first-api-go/router"
)

var setting config.SettingProvider

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

	redisClient := redis.NewClient()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := controller.NewUserHandler(userService)

	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository, userService)
	postHandler := controller.NewPostHandler(postService)

	likeRepository := repository.NewLikeRepository(db)
	likeService := service.NewLikeService(likeRepository, userService, postService, redisClient)
	likeHandler := controller.NewLikeHandler(likeService)

	replyRepository := repository.NewReplyRepository(db)
	replyService := service.NewReplyService(replyRepository, userService, postService, redisClient)
	replyHandler := controller.NewReplyHandler(replyService)

	conn, channel := config.InitBroker()

	defer channel.Close()
	defer conn.Close()

	exchange := "topic_report"
	routingKey := "post.report.created"

	reportRepository := repository.NewReportRepository(db)
	reportProducer := queue.NewReportProducer(channel, exchange, routingKey)
	reportService := reportService.NewReportService(reportRepository, reportProducer, postService, userService)
	reportHandler := controller.NewReportHandler(reportService)

	r := router.SetupRouter(userHandler, postHandler, likeHandler, replyHandler, reportHandler)

	port := setting.GetServer().Port

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Falha ao iniciar servidor:", err)
	}
}
