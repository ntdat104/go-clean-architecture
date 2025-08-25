package repo

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConfig struct {
	Driver                  string
	Url                     string
	ConnMaxLifetimeInMinute int
	MaxOpenConns            int
	MaxIdleConns            int
}

func NewDB(conf DatabaseConfig) (*sql.DB, error) {
	db, err := newDatabase(conf)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newDatabase(conf DatabaseConfig) (*sql.DB, error) {
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

	return db, err
}
