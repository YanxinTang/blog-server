package server

import (
	"log"

	"github.com/YanxinTang/blog-server/config"
	"github.com/YanxinTang/blog-server/middleware"
	"github.com/YanxinTang/blog-server/model"
	"github.com/YanxinTang/blog-server/router"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	conf, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
	pool, err := config.GetDBConnectionPool(conf.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	store, err := config.GetCookieStore(*conf)
	if err != nil {
		log.Fatal(err)
	}
	model.Setup(pool)

	svr := gin.Default()
	svr.Use(sessions.Sessions("sessionid", store), middleware.ErrorHandler())
	router.SetupRouter(svr)

	return svr
}
