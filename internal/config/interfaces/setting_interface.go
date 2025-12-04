package interfaces

import (
	types "github.com/EmersonRabelo/first-api-go/internal/config/types"
)

type SettingProvider interface {
	GetEnvironment() string
	GetServer() types.Server
	GetDatabase() types.Database
	IsProd() bool
	IsTest() bool
	IsLocal() bool
}
