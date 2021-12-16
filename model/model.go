package model

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool
var ctx context.Context

type BaseModel struct {
	ID        uint64    `json:"id" db:"id"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

func Setup(pool *pgxpool.Pool) {
	db = pool
	ctx = context.Background()
}
