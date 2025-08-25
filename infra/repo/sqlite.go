package repo

import (
	"database/sql"
	"io/ioutil"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteConfig struct {
	// Path to your SQLite database file.
	// Example: "./app.db" or ":memory:" for in-memory DB
	Path                  string
	SchemaPath            string // path to schema.sql
	ConnMaxLifetimeMinute int
	MaxOpenConns          int
	MaxIdleConns          int
}

func NewSQLiteDB(conf SQLiteConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", conf.Path)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(conf.ConnMaxLifetimeMinute))
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Run schema migration
	if err := runSchema(db, conf.SchemaPath); err != nil {
		return nil, err
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
