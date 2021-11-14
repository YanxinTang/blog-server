package model

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/YanxinTang/blog/server/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool
var ctx context.Context

type BaseModel struct {
	ID        uint64    `json:"id" db:"id"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type qLogger struct {
}

func (q *qLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	fmt.Printf("SQL: %s\nARGS: %v\n", data["sql"], data["args"])
}

func init() {
	database := &config.Config.Database

	ctx = context.Background()
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		database.User,
		database.Password,
		database.Host,
		database.Port,
		database.DBName,
	)

	var err error
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	config.ConnConfig.Logger = &qLogger{}

	db, err = pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
}
