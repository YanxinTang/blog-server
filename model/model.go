package model

import (
	"fmt"
	"time"

	"github.com/YanxinTang/blog/server/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

type BaseModel struct {
	ID        uint64    `json:"id" db:"id"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

func init() {
	mysql := &config.Config.Mysql
	connect := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true",
		mysql.User,
		mysql.Password,
		mysql.Host,
		mysql.Port,
		mysql.Database,
	)
	DB = sqlx.MustConnect("mysql", connect)
}
