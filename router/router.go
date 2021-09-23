package router

import (
	"encoding/gob"

	"github.com/YanxinTang/blog/server/controller"
	"github.com/YanxinTang/blog/server/middleware"
	"github.com/YanxinTang/blog/server/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func init() {
	gob.Register(&model.User{})
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionid", store), middleware.ErrorHandler())

	api := r.Group("api")
	public := api.Group("")
	{
		public.POST("signup", controller.Signup)
		public.POST("signin", controller.Signin)
		public.POST("signout", controller.Signout)
		public.GET("login/session", controller.GetLoginSession)

		public.GET("articles", controller.GetArticles)
		public.GET("articles/:articleID", controller.GetArticle)
		public.GET("articles/:articleID/comments", controller.GetArticleComments)
		public.POST("articles/:articleID/comments", controller.CreateComment)
		public.GET("categories", controller.GetCategories)

		public.GET("setting", controller.GetSetting)
	}

	protected := api.Group("admin", middleware.Auth())
	{
		protected.GET("overview", controller.Overview) // 总览

		// 分类
		protected.POST("/categories", controller.CareteCategory)
		protected.PUT("/categories/:categoryID", controller.UpdateCategory)
		protected.DELETE("/categories/:categoryID", controller.DeleteCategory)

		// 文章
		protected.POST("articles", controller.CreateArticle)
		protected.DELETE("/articles/:articleID", controller.DeleteArticle)
		protected.PUT("/articles/:articleID", controller.UpdateArticle)
		protected.DELETE("/articles/:articleID/comment/:commentID", controller.DeleteComment)

		// 草稿
		protected.GET("drafts", controller.GetDrafts)
		protected.POST("drafts", controller.CreateDraft)
		protected.POST("drafts/:draftID", controller.PublishDraft)
		protected.GET("drafts/:draftID", controller.GetDraft)
		protected.PUT("drafts/:draftID", controller.UpdateDraft)
		protected.DELETE("drafts/:draftID", controller.DeleteDraft)

		// 工具接口
		public.POST("setting", controller.SetSetting)
	}

	return r
}
