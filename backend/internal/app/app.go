package app

import (
	"backend/global"
	roleWire "backend/internal/infrastructures/wire/role"

	// userWire "backend/internal/infrastructures/wire/user"
	"backend/internal/initializations"
	// "backend/internal/presentations/http"
	// "backend/internal/presentations/http/v1/router"
)

func Run() {
	// ===== config =====
	initializations.LoadConfig()

	// ===== logger =====
	initializations.InitLogger()
	global.Logger.Info("Config log successfully")

	// ===== postgres =====
	pgDb := initializations.NewPostgres()
	global.Logger.Info("Initializing Postgres successfully")

	// ===== service =====
	roleService := roleWire.NewRoleService(pgDb)

	// init role cache
	go initializations.NewRolesCache(roleService)

	// ===== router =====
	router := initializations.InitRouter(pgDb)
	
	// ===== server =====
	server := initializations.NewServer(router)
	initializations.RunServer(server)
}
