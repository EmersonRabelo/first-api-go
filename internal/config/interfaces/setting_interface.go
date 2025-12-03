package interfaces

import (
	model "github.com/EmersonRabelo/first-api-go/internal/config/model"
)

type SettingProvider interface {
	GetEnvironment() string
	GetServer() model.Server
	GetDatabase() model.Database
	IsProd() bool
	IsTest() bool
	IsLocal() bool
}
