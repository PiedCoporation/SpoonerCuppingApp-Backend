package app

import (
	"backend/config"
	"backend/internal/controller/http"
	"backend/internal/controller/http/v1/router"
	roleWire "backend/internal/infrastructure/wire/role"
	userWire "backend/internal/infrastructure/wire/user"
	"backend/internal/initialization"
	"backend/pkg/logger"
)

func Run(cfg *config.Config) {
	// ===== logger =====
	l := logger.New(cfg.Logger)
	l.Info("Config log successfully")

	// ===== postgres =====
	pgDb := initialization.NewPostgres(&cfg.Postgres, l)
	l.Info("Initializing Postgres successfully")

	// ===== service =====
	// user
	userAuthSerice := userWire.NewUserAuthService(cfg, pgDb)
	// role
	roleService := roleWire.NewRoleService(pgDb)

	// init role cache
	go initialization.NewRolesCache(roleService, l)

	// ===== router =====
	userServiceSet := &router.UserServiceSet{
		UserAuthService: userAuthSerice,
	}
	router := http.NewRouter(cfg, userServiceSet)

	// ===== server =====
	server := initialization.NewServer(&cfg.HTTP, router)
	initialization.RunServer(server, &cfg.HTTP, l)
}
