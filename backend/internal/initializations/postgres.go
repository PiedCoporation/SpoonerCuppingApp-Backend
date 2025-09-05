package initializations

import (
	"backend/config"
	"backend/global"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres() *gorm.DB {
	pgCfg := global.Config.Postgres

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		pgCfg.Host, pgCfg.Username, pgCfg.Password, pgCfg.Dbname, pgCfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		global.Logger.Fatal("NewPostgres initializations error", zap.Error(err))
	}

	sqlDb, err := db.DB()
	if err != nil {
		global.Logger.Fatal("Get sql.DB error", zap.Error(err))
	}
	if err := sqlDb.Ping(); err != nil {
		global.Logger.Fatal("Database not reachable", zap.Error(err))
	}

	SetPool(sqlDb, &pgCfg)

	return db
}

func SetPool(sqlDb *sql.DB, pgCfg *config.Postgres) {
	sqlDb.SetMaxIdleConns(pgCfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(pgCfg.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(pgCfg.ConnMaxLifetime * time.Second)
}
