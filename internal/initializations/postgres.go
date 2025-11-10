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
    // Align with golang-migrate schema_migrations table (single-row with version/dirty)
    if _, err := sqlDb.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version bigint not null, dirty boolean not null)`); err != nil {
        return err
    }
    // Ensure a single row exists
    var count int
    if err := sqlDb.QueryRow(`SELECT COUNT(1) FROM schema_migrations`).Scan(&count); err != nil {
        return err
    }
    if count == 0 {
        if _, err := sqlDb.Exec(`INSERT INTO schema_migrations (version, dirty) VALUES (0, false)`); err != nil {
            return err
        }
    }
    // Read current version and dirty flag
    var currentVersion int64
    var dirty bool
    if err := sqlDb.QueryRow(`SELECT version, dirty FROM schema_migrations LIMIT 1`).Scan(&currentVersion, &dirty); err != nil {
        return err
    }
    if dirty {
        return fmt.Errorf("previous migration left schema_migrations.dirty = true; manual intervention required")
    }

    // Discover .up.sql files
    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    type mig struct{ version int64; name string; path string }
    list := make([]mig, 0, len(entries))
    for _, e := range entries {
        if e.IsDir() { continue }
        name := e.Name()
        if !strings.HasSuffix(name, ".up.sql") { continue }
        // parse leading numeric version before first underscore
        verStr := name
        if idx := strings.IndexByte(name, '_'); idx > 0 {
            verStr = name[:idx]
        }
        var v int64
        // tolerate leading zeros
        fmt.Sscanf(verStr, "%d", &v)
        list = append(list, mig{version: v, name: name, path: dir + "/" + name})
    }
    sort.Slice(list, func(i, j int) bool { return list[i].version < list[j].version })

    // Apply pending migrations strictly by version
    for _, m := range list {
        if m.version <= currentVersion {
            continue
        }
        content, err := os.ReadFile(m.path)
        if err != nil { return err }
        // mark dirty
        if _, err := sqlDb.Exec(`UPDATE schema_migrations SET dirty = true`); err != nil { return err }
        // apply
        if _, err := sqlDb.Exec(string(content)); err != nil {
            // keep dirty=true so operator can see failure
            return fmt.Errorf("migration %s failed: %w", m.name, err)
        }
        // set new version and clear dirty
        if _, err := sqlDb.Exec(`UPDATE schema_migrations SET version = $1, dirty = false`, m.version); err != nil { return err }
        currentVersion = m.version
        global.Logger.Info("applied migration", zap.String("file", m.name))
    }

    return nil
}