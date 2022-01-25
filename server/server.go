package server

import (
	"github.com/YanxinTang/blog-server/config"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/middleware"
	"github.com/YanxinTang/blog-server/model"
	"github.com/YanxinTang/blog-server/router"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func New() *gin.Engine {
	conf, err := config.ParseConfig()
	if err != nil {
		log.Fatal("failed to parse config", zap.Error(err))
	}
	pool, err := config.GetDBConnectionPool(conf.Postgres)
	if err != nil {
		log.Fatal("failed to connect to db", zap.Error(err))
	}
	store, err := config.GetCookieStore(*conf)
	if err != nil {
		log.Fatal("failed to get cookie store", zap.Error(err))
	}
	model.Setup(pool)

	svr := gin.Default()
	svr.Use(sessions.Sessions("sessionid", store), middleware.ErrorHandler())
	router.SetupRouter(svr)

	return svr
}
