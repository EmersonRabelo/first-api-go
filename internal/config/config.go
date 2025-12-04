package config

import (
	"github.com/EmersonRabelo/first-api-go/internal/config/interfaces"
	types "github.com/EmersonRabelo/first-api-go/internal/config/types"
)

var AppSetting interfaces.SettingProvider

func GetSetting() interfaces.SettingProvider {
	AppSetting = types.LoadSetting()

	return AppSetting
}
