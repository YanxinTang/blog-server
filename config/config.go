package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"gopkg.in/yaml.v2"
)

var SigninKey = []byte("blog")
var CaptchaExpiration = time.Minute * 5

// PostgresConfig persists the config for our PostgreSQL database connection
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type CookieStoreConfig struct {
	Database string `yaml:"database"`
	Secret   string `yaml:"secret"`
}

type Config struct {
	Postgres    PostgresConfig    `yaml:"postgres"`
	CookieStore CookieStoreConfig `yaml:"cookiestore"`
}

func ParseConfig() (*Config, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	configFilePath := filepath.Join(filepath.Dir(ex), "./conf/config.yaml")
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.NewDecoder(file).Decode(&config)
	return &config, err
}

func GetEntClient(conf *Config) (*ent.Client, error) {
	connection := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.User,
		conf.Postgres.Database,
		conf.Postgres.Password,
	)
	return ent.Open("postgres", connection)
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

func GetCookieStore(conf *Config) (cookie.Store, error) {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.CookieStore.Database,
	)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	store, err := postgres.NewStore(db, []byte(conf.CookieStore.Secret))
	if err != nil {
		return nil, err
	}
	return store, nil
}

func GetDefaultConnectionPool() (*pgxpool.Pool, error) {
	configuration, err := ParseConfig()
	if err != nil {
		return nil, err
	}
	return GetDBConnectionPool(configuration.Postgres)
}
