package model

import (
	"context"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool
var DB *pgxpool.Pool
var ctx context.Context

type BaseModel struct {
	ID        uint64    `json:"id" db:"id"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
}

func Setup(pool *pgxpool.Pool) {
	db = pool
	DB = pool
	ctx = context.Background()
}
