package config

import (
	"github.com/EmersonRabelo/first-api-go/internal/config/interfaces"
	"github.com/EmersonRabelo/first-api-go/internal/config/model"
)

var AppSetting interfaces.SettingProvider

func GetSetting() interfaces.SettingProvider {
	AppSetting = model.LoadSetting()

	return AppSetting
}
