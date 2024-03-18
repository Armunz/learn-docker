package config

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

type MySQLConfig struct {
	dsn                   string
	maxOpenConn           int
	maxConnLifeTimeSecond int
	maxIdleConn           int
}

func NewMySQL(dsn string, maxOpenConn int, maxConnLifeTimeSecond int, maxIdleConn int) MySQLConfig {
	return MySQLConfig{
		dsn:                   dsn,
		maxOpenConn:           maxOpenConn,
		maxConnLifeTimeSecond: maxConnLifeTimeSecond,
		maxIdleConn:           maxIdleConn,
	}
}

func (m MySQLConfig) Connect() *sql.DB {
	db, err := sql.Open("mysql", m.dsn)
	if err != nil {
		log.Panic().Err(err).Caller().Msg("failed to open mysql database")
	}

	db.SetMaxOpenConns(m.maxOpenConn)
	db.SetConnMaxLifetime(time.Duration(m.maxConnLifeTimeSecond) * time.Second)
	db.SetMaxIdleConns(m.maxIdleConn)

	if err := db.Ping(); err != nil {
		log.Panic().Err(err).Caller().Msg("failed to ping mysql database")
	}

	return db
}
