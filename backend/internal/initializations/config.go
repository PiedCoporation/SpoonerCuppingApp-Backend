package initializations

import (
	"backend/config"
	"backend/global"

	"github.com/caarlos0/env/v11"
)

func LoadConfig() {
	global.Config = &config.Config{}
	if err := env.Parse(global.Config); err != nil {
		panic(err)
	}
}
