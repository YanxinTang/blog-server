package router

import (
	"encoding/gob"

	"github.com/YanxinTang/blog-server/internal/app/controller"
	"github.com/YanxinTang/blog-server/internal/app/middleware"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/gin-gonic/gin"
)

func init() {
	gob.Register(&model.User{})
}

func SetupRouter(server *gin.Engine) {
	api := server.Group("api")
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
		public.GET("categories/:categoryID/articles", controller.GetCategoryArticles)

		public.GET("setting", controller.GetSetting)

		public.GET("storages/:storageID/:key", controller.GetStorageObject)

		public.GET("captcha", controller.GetCapacha)
	}

	protected := api.Group("admin", middleware.Auth())
	{
		protected.GET("overview", controller.Overview) // 总览
		protected.GET("overview/storage", controller.OverviewStorage)

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

		// 存储
		protected.GET("storages", controller.GetStorages)
		protected.GET("storages/:storageID", controller.GetStorage)
		protected.POST("storages", controller.CreateStorage)
		protected.PUT("storages/:storageID", controller.UpdateStorage)
		protected.DELETE("storages/:storageID", controller.DeleteStorage)
		protected.GET("storages/:storageID/objects", controller.GetStorageObjects)
		protected.DELETE("storages/:storageID/object", controller.DeleteStorageObject)
		protected.PUT("storages/:storageID/upload", controller.PutStorageObject)

		// 工具接口
		protected.POST("setting", controller.SetSetting)
	}
}