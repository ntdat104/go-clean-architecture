package repo

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseConfig struct {
	Driver                  string // "mysql" or "sqlite3"
	Url                     string // DSN (e.g. MySQL: "user:pass@tcp(localhost:3306)/dbname"; SQLite: "./app.db" or ":memory:")
	SchemaPath              string // path to schema.sql (optional)
	ConnMaxLifetimeInMinute int
	MaxOpenConns            int
	MaxIdleConns            int
}

func NewDB(conf DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(conf.Driver, conf.Url)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(conf.ConnMaxLifetimeInMinute))
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Run schema migration if SchemaPath is provided
	if conf.SchemaPath != "" {
		if err := runSchema(db, conf.SchemaPath); err != nil {
			return nil, fmt.Errorf("failed to run schema: %w", err)
		}
	}

	return db, nil
}

func runSchema(db *sql.DB, schemaPath string) error {
	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(schema))
	return err
}
