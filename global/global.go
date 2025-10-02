package global

import (
	"backend/config"
	"backend/pkg/logger"
)

var (
	Config *config.Config
	Logger logger.Interface
)
