package initializations

import (
	"backend/global"
	"backend/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.New(global.Config.Logger)
}
