package server

import (
	"github.com/YanxinTang/blog-server/internal/app/middleware"
	"github.com/YanxinTang/blog-server/internal/app/router"
	"github.com/YanxinTang/blog-server/internal/app/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func New(sessionStore sessions.Store) *gin.Engine {
	go service.CaptchaStoreGC()

	svr := gin.Default()
	svr.Use(sessions.Sessions("sessionid", sessionStore), middleware.ErrorHandler())
	router.SetupRouter(svr)

	return svr
}
