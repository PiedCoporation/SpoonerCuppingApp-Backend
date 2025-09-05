package app

import (
	"backend/global"
	"backend/internal/controller/http"
	"backend/internal/controller/http/v1/router"
	roleWire "backend/internal/infrastructure/wire/role"
	userWire "backend/internal/infrastructure/wire/user"
	"backend/internal/initialization"
)

func Run() {
	// ===== config =====
	initialization.LoadConfig()

	// ===== logger =====
	initialization.InitLogger()
	global.Logger.Info("Config log successfully")

	// ===== postgres =====
	pgDb := initialization.NewPostgres()
	global.Logger.Info("Initializing Postgres successfully")

	// ===== service =====
	// user
	userAuthSerice := userWire.NewUserAuthService(pgDb)
	// role
	roleService := roleWire.NewRoleService(pgDb)

	// init role cache
	go initialization.NewRolesCache(roleService)

	// ===== router =====
	userServiceSet := &router.UserServiceSet{
		UserAuthService: userAuthSerice,
	}
	router := http.NewRouter(userServiceSet)

	// ===== server =====
	server := initialization.NewServer(router)
	initialization.RunServer(server)
}
