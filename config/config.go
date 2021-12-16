package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
)

var SigninKey = []byte("blog")

// PostgresConfig persists the config for our PostgreSQL database connection
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Config struct {
	Postgres PostgresConfig `json:"postgres"`
}

func ParseConfig() (*Config, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	configFilePath := filepath.Join(filepath.Dir(ex), "./config/config.json")
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.NewDecoder(file).Decode(&config)
	return &config, err
}

func GetDBConnectionPool(postgreConfig PostgresConfig) (*pgxpool.Pool, error) {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		postgreConfig.User,
		postgreConfig.Password,
		postgreConfig.Host,
		postgreConfig.Port,
		postgreConfig.Database,
	)

	var err error
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}
	return pgxpool.ConnectConfig(context.Background(), config)
}

func GetDBMigrate(pool *pgxpool.Pool) (*migrate.Migrate, error) {
	db := stdlib.OpenDB(*pool.Config().ConnConfig)
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance("file://migrations", "pgx", driver)
}

func GetDefaultConnectionPool() (*pgxpool.Pool, error) {
	configuration, err := ParseConfig()
	if err != nil {
		return nil, err
	}
	return GetDBConnectionPool(configuration.Postgres)
}

func GetDefaultMigrate() (*migrate.Migrate, error) {
	pool, err := GetDefaultConnectionPool()
	if err != nil {
		return nil, err
	}
	return GetDBMigrate(pool)
}
