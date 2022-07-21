package controller

import (
	"net/http"

	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/app/service"
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/YanxinTang/blog-server/internal/pkg/page"
	"github.com/YanxinTang/blog-server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateArticleReqBody struct {
	Title      string              `json:"title"  binding:"required"`
	CategoryID int                 `json:"categoryID" binding:"required"`
	Content    string              `json:"content" binding:"required"`
	Status     model.ArticleStatus `json:"status" binding:"required"`
}

func CreateArticle(c *gin.Context) {
	var createArticleReqBody CreateArticleReqBody
	if err := c.Bind(&createArticleReqBody); err != nil {
		return
	}

	cai := model.CreateArticleInput{
		Title:      createArticleReqBody.Title,
		CategoryID: createArticleReqBody.CategoryID,
		Content:    createArticleReqBody.Content,
		Status:     createArticleReqBody.Status,
	}

	article, err := model.CreateArticle(common.Context, common.Client)(cai)
	if err != nil {
		log.Warn("failing creating article", zap.Error(err))
		c.Error(e.ERROR_INTERVAL_ERROR)
		return
	}
	c.JSON(http.StatusOK, article)
}

type UpdateArticleReqBody struct {
	CategoryID int                 `json:"categoryID"`
	Title      string              `json:"title"`
	Content    string              `json:"content"`
	Status     model.ArticleStatus `json:"status"`
}

func UpdateArticle(c *gin.Context) {
	articleID, err := utils.GetID(c, "articleID")
	if err != nil {
		c.Error(err)
		return
	}

	var updateArticleReqBody UpdateArticleReqBody
	if err := c.BindJSON(&updateArticleReqBody); err != nil {
		return
	}

	uai := model.UpdateArticleInput{
		ID:         articleID,
		CategoryID: updateArticleReqBody.CategoryID,
		Title:      updateArticleReqBody.Title,
		Content:    updateArticleReqBody.Content,
		Status:     updateArticleReqBody.Status,
	}

	article, err := model.UpdateArticle(common.Context, common.Client)(uai)
	if err != nil {
		log.Warn("failing updating article", zap.Error(err))
		c.Error(e.ERROR_BAD_REQUEST)
		return
	}

	c.JSON(http.StatusOK, article)
}

func DeleteArticle(c *gin.Context) {
	articleID, err := utils.GetID(c, "articleID")
	if err != nil {
		return
	}

	if err := service.DeleteArticle(articleID); err != nil {
		c.Error(e.ERROR_BAD_REQUEST)
		return
	}
	c.Status(http.StatusOK)
}

func GetArticle(c *gin.Context) {
	articleID, err := utils.GetID(c, "articleID")
	if err != nil {
		c.Error(e.ERROR_BAD_REQUEST)
		return
	}

	article, err := model.GetArticle(common.Context, common.Client)(articleID, model.StatusNil)
	if err != nil {
		c.Error(e.ERROR_RESOURCE_NOT_FOUND)
		return
	}

	c.JSON(http.StatusOK, article)
}

type GetArticlesReqQuery struct {
	CategoryID int                 `form:"categoryID"`
	Status     model.ArticleStatus `form:"status,default=0"`
}

// GetArticles 获取所有的文章
func GetArticles(c *gin.Context) {
	var getArticlesReqQuery GetArticlesReqQuery
	if err := c.BindQuery(&getArticlesReqQuery); err != nil {
		c.Error(err)
		return
	}
	pagination := page.NewPagination()
	if err := c.BindQuery(&pagination); err != nil {
		return
	}

	articles, pagination, err := service.GetArticles(model.CategoryNil, model.StatusNil, pagination)
	if err != nil {
		log.Warn(
			"failing getting articles",
			zap.Int("categroyID", getArticlesReqQuery.CategoryID),
			zap.Int8("status", int8(getArticlesReqQuery.Status)),
			zap.Error(err),
		)
		c.Error(e.ERROR_RESOURCE_NOT_FOUND)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":   articles,
		"pagination": pagination,
	})
}

func PublicGetArticle(c *gin.Context) {
	articleID, err := utils.GetID(c, "articleID")
	if err != nil {
		c.Error(e.ERROR_BAD_REQUEST)
		return
	}

	article, err := service.GetPublishedArticle(articleID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, article)
}

// GetCategoryArticles 获取某个分类下的所有文章
func PublicGetCategoryArticles(c *gin.Context) {
	categoryID, err := utils.GetID(c, "categoryID")
	if err != nil {
		c.Error(e.ERROR_BAD_REQUEST)
		return
	}

	pagination := page.NewPagination()
	if err := c.BindQuery(&pagination); err != nil {
		return
	}

	articles, pagination, err := service.GetCategoryPublishedArticles(categoryID, pagination)
	if err != nil {
		c.Error(err)
		return
	}

	if !pagination.IsValid() {
		c.Error(e.ERROR_RESOURCE_NOT_FOUND)
		return
	}

	for i := range articles {
		articles[i].Content = utils.Summary(articles[i].Content)
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":   articles,
		"pagination": pagination,
	})
}

// PublicGetArticles 获取已发布的所有的文章
func PublicGetArticles(c *gin.Context) {
	pagination := page.NewPagination()
	if err := c.BindQuery(&pagination); err != nil {
		return
	}

	articles, pagination, err := service.GetPublishedArticles(pagination)
	if err != nil {
		c.Error(e.ERROR_RESOURCE_NOT_FOUND)
		return
	}

	if !pagination.IsValid() {
		c.Error(e.ERROR_RESOURCE_NOT_FOUND)
		return
	}

	for i := range articles {
		articles[i].Content = utils.Summary(articles[i].Content)
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":   articles,
		"pagination": pagination,
	})
}
