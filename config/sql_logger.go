package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type sqlLogger struct {
}

func (l *sqlLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	fmt.Printf("SQL: %s\nARGS: %v\n", data["sql"], data["args"])
}
