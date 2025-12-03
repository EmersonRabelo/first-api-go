package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	gormConfig := &gorm.Config{}

	database := AppSetting.GetDatabase()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		database.Host,
		database.Port,
		database.User,
		database.Pwd,
		database.Name,
		database.SSLMode,
	)

	if AppSetting.IsProd() {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	} else {
		gormConfig.Logger = logger.New(
			log.New(os.Stdin, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
	}
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)

	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Falha ao configurar connection pool:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Conexão com banco de dados estabelecida com sucesso!")
}

func GetDB() *gorm.DB {
	return DB
}
