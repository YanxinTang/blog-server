package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/YanxinTang/blog-server/e"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/model"
	"github.com/YanxinTang/blog-server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateCommentReq struct {
	VerifyCaptchaReq
	Username string `json:"username"`
	Content  string `json:"content" binding:"required"`
}

func CreateComment(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	var createCommentReq CreateCommentReq
	if err := c.BindJSON(&createCommentReq); err != nil {
		log.Warn("create comment binding error", zap.Error(err))
		return
	}

	if apierr := service.VerifyCaptcha(createCommentReq.Key, createCommentReq.Text); apierr != nil {
		c.Error(apierr)
		return
	}

	comment := model.Comment{
		ArticleID: articleID,
		Username:  strings.TrimSpace(createCommentReq.Username),
		Content:   createCommentReq.Content,
	}

	if comment.Username == "" {
		comment.Username = "匿名"
	}

	comment, err = model.CreateComment(comment)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, comment)
}

func GetArticleComments(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	comments, err := model.GetArticleComments(articleID)
	if err != nil {
		c.Error(e.ERROR_RESOURCE_NOT_FOUND)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}
