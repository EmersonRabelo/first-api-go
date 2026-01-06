package main

import (
	"fmt"
	"log"

	config "github.com/EmersonRabelo/first-api-go/internal/config"
	controller "github.com/EmersonRabelo/first-api-go/internal/controller"
	database "github.com/EmersonRabelo/first-api-go/internal/database"
	"github.com/EmersonRabelo/first-api-go/internal/handler"
	"github.com/EmersonRabelo/first-api-go/internal/queue"
	redis "github.com/EmersonRabelo/first-api-go/internal/redis"
	repository "github.com/EmersonRabelo/first-api-go/internal/repository"
	service "github.com/EmersonRabelo/first-api-go/internal/service"
	consumer "github.com/EmersonRabelo/first-api-go/internal/service/consumer"
	reportService "github.com/EmersonRabelo/first-api-go/internal/service/report"
	router "github.com/EmersonRabelo/first-api-go/router"
)

var setting config.SettingProvider

func init() {
	fmt.Println("Application initializing...")

	setting = config.GetSetting()

	config.InitDatabase()
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
	routingKeyProducer := "post.report.created"
	routingKeyConsumer := "post.report.response"
	queueName := "q.report.response"

	reportRepository := repository.NewReportRepository(db)
	reportProducer := queue.NewReportProducer(channel, exchange, routingKeyProducer)
	reportService := reportService.NewReportService(reportRepository, reportProducer, postService, userService)
	reportHandler := controller.NewReportHandler(reportService)

	reportConsumerService := consumer.NewConsumerReportService(reportRepository)
	handler := handler.NewReportHandler(reportConsumerService)
	reportConsumer := queue.NewReportConsumer(channel, exchange, routingKeyConsumer, queueName, handler)

	go func() {
		if err := reportConsumer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	r := router.SetupRouter(userHandler, postHandler, likeHandler, replyHandler, reportHandler)

	port := setting.GetServer().Port

	fmt.Println("Port", port)

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Falha ao iniciar servidor:", err)
	}

	fmt.Println("Initialized.")
}
