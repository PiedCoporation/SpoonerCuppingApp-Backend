package initializations

import (
	"backend/config"
	"backend/global"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres() *gorm.DB {
	pgCfg := global.Config.Postgres

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
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

    // Run SQL migrations on startup using simple file runner
    if err := RunMigrations(sqlDb, "migrations"); err != nil {
        global.Logger.Warn("migrations failed", zap.Error(err))
    }

	return db
}

func SetPool(sqlDb *sql.DB, pgCfg *config.Postgres) {
	sqlDb.SetMaxIdleConns(pgCfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(pgCfg.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(pgCfg.ConnMaxLifetime * time.Second)
}

// RunMigrations executes .up.sql files in lexicographical order.
// It is a minimal runner that applies each migration exactly once by tracking filenames in a table.
func RunMigrations(sqlDb *sql.DB, dir string) error {
    // Ensure migrations table exists
    if _, err := sqlDb.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (id SERIAL PRIMARY KEY, filename TEXT UNIQUE NOT NULL, applied_at TIMESTAMP NOT NULL DEFAULT NOW())`); err != nil {
        return err
    }

    // Read applied migrations
    applied := map[string]struct{}{}
    rows, err := sqlDb.Query(`SELECT filename FROM schema_migrations`)
    if err != nil {
        return err
    }
    defer rows.Close()
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            return err
        }
        applied[name] = struct{}{}
    }
    if err := rows.Err(); err != nil {
        return err
    }

    // Discover .up.sql files
    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    type mig struct{ name string; path string }
    list := make([]mig, 0, len(entries))
    for _, e := range entries {
        if e.IsDir() { continue }
        name := e.Name()
        if !strings.HasSuffix(name, ".up.sql") { continue }
        list = append(list, mig{name: name, path: dir + "/" + name})
    }
    sort.Slice(list, func(i, j int) bool { return list[i].name < list[j].name })

    // Apply pending migrations
    for _, m := range list {
        if _, ok := applied[m.name]; ok { continue }
        content, err := os.ReadFile(m.path)
        if err != nil { return err }
        if _, err := sqlDb.Exec(string(content)); err != nil { return fmt.Errorf("migration %s failed: %w", m.name, err) }
        if _, err := sqlDb.Exec(`INSERT INTO schema_migrations (filename) VALUES ($1)`, m.name); err != nil { return err }
        global.Logger.Info("applied migration", zap.String("file", m.name))
    }

    return nil
}