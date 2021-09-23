package main

import (
	"github.com/YanxinTang/blog/server/router"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	engine := router.SetupRouter()

	// Listen and Server in 0.0.0.0:8080
	engine.Run(":8000")
}
