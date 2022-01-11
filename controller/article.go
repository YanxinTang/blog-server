package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/YanxinTang/blog-server/e"
	"github.com/YanxinTang/blog-server/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const PerPage uint64 = 10

type apiGetArticlesQueryModel struct {
	Page uint64 `form:"page,default=1"`
}

// GetArticles 获取所有的文章
func GetArticles(c *gin.Context) {
	var apiGetArticlesQuery apiGetArticlesQueryModel
	if err := c.ShouldBindQuery(&apiGetArticlesQuery); err != nil {
		c.Error(err)
		return
	}
	pagination := model.NewPagination()
	if err := c.Bind(&pagination); err != nil {
		return
	}
	var err error
	pagination.Total, err = model.ArticlesCount()
	if err != nil {
		c.Error(err)
		return
	}
	if !pagination.IsValid() {
		c.Error(e.ERROR_RESOURCE_NOT_FOUND)
		return
	}

	articles, err := model.GetPublishedArticles(pagination)
	if err != nil {
		c.Error(e.ERROR_RESOURCE_NOT_FOUND)
		return
	}

	for i := range articles {
		articles[i].Content = articles[i].Summary()
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":   articles,
		"pagination": pagination,
	})
}

// GetCategoryArticles 获取某个分类下的所有文章
func GetCategoryArticles(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("categoryID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	pagination := model.NewPagination()
	if err := c.BindQuery(&pagination); err != nil {
		return
	}

	pagination.Total, err = model.CategoryArticlesCount(categoryID)
	if err != nil {
		c.Error(err)
		return
	}

	if !pagination.IsValid() {
		c.Error(e.ERROR_RESOURCE_NOT_FOUND)
		return
	}

	articles, err := model.GetCategoryPublishedArticles(categoryID, pagination)
	if err != nil {
		c.Error(err)
		return
	}

	for i := range articles {
		articles[i].Content = articles[i].Summary()
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":   articles,
		"pagination": pagination,
	})
}

func GetArticle(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	article, err := model.GetPublishedArticle(articleID)
	if err != nil {
		c.Error(e.ERROR_RESOURCE_NOT_FOUND)
		return
	}
	c.JSON(http.StatusOK, article)
}

type apiCreateArticleModel struct {
	Title      string `json:"title"  binding:"required"`
	CategoryID uint64 `json:"categoryID" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

func CreateArticle(c *gin.Context) {
	var apiCreateArticle apiCreateArticleModel
	if err := c.Bind(&apiCreateArticle); err != nil {
		return
	}

	session := sessions.Default(c)
	userID := session.Get("userID").(uint64)

	article := model.Article{
		Title:      apiCreateArticle.Title,
		CategoryID: apiCreateArticle.CategoryID,
		Content:    apiCreateArticle.Content,
	}

	lastInsertArticle, err := model.CreateArticle(userID, article)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lastInsertArticle)
}

func UpdateArticle(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	var article model.Article
	if err := c.ShouldBind(&article); err != nil {
		c.Error(err)
		return
	}

	article.ID = articleID

	if _, err := model.UpdateArticle(article); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func DeleteArticle(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	if err := model.DeleteArticle(articleID); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

/* Draft */

func GetDrafts(c *gin.Context) {
	pagination := model.NewPagination()
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.Error(err)
		return
	}

	drafts, err := model.GetDrafts(pagination)
	if err != nil {
		c.Error(err)
		return
	}
	pagination.Total, err = model.ArticlesCount()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"drafts":     drafts,
		"pagination": pagination,
	})
}

func GetDraft(c *gin.Context) {
	draftID, err := strconv.ParseUint(c.Param("draftID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	draft, err := model.GetDraft(draftID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, draft)
}

func CreateDraft(c *gin.Context) {
	var article model.Article
	if err := c.BindJSON(&article); err != nil {
		c.Error(err)
		return
	}

	session := sessions.Default(c)
	userID := session.Get("userID").(uint64)
	article.Status = model.StatusDraft

	lastInsertDraft, err := model.CreateArticle(userID, article)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lastInsertDraft)
}

func UpdateDraft(c *gin.Context) {
	draftID, err := strconv.ParseUint(c.Param("draftID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	var article model.Article
	if err := c.ShouldBind(&article); err != nil {
		c.Error(err)
		return
	}

	article.ID = draftID

	if _, err := model.UpdateArticle(article); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func PublishDraft(c *gin.Context) {
	draftID, err := strconv.ParseUint(c.Param("draftID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	if err := model.PublishDraft(draftID); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func DeleteDraft(c *gin.Context) {
	draftID, err := strconv.ParseUint(c.Param("draftID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	if err := model.DeleteArticle(draftID); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func DeleteComment(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	commentID, err := strconv.ParseUint(c.Param("commentID"), 10, 64)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	if _, err = model.DeleteComment(commentID); err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	session := sessions.Default(c)
	session.AddFlash("删除成功", "successMsgs")
	session.Save()
	c.Redirect(http.StatusFound, fmt.Sprintf("/articles/%d", articleID))
}
