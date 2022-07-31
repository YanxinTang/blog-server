package controller

import (
	"net/http"
	"strings"

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

type CreateCommentReqBody struct {
	VerifyCaptchaReqBody
	Username string `json:"username"`
	Content  string `json:"content" binding:"required"`
}

func CreateComment(c *gin.Context) {
	articleID, err := utils.GetID(c, "articleID")
	if err != nil {
		return
	}
	var createCommentReqBody CreateCommentReqBody
	if err := c.BindJSON(&createCommentReqBody); err != nil {
		log.Warn("create comment binding error", zap.Error(err))
		return
	}

	if err := service.VerifyCaptcha(createCommentReqBody.Key, createCommentReqBody.Text); err != nil {
		c.Error(err)
		return
	}

	cci := model.CreateCommentInput{
		ArticleID: articleID,
		Username:  strings.TrimSpace(createCommentReqBody.Username),
		Content:   createCommentReqBody.Content,
	}

	if cci.Username == "" {
		cci.Username = "匿名"
	}

	comment, err := model.CreateComment(common.Context, common.Client)(cci)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	commentID, err := utils.GetID(c, "commentID")
	if err != nil {
		c.Error(e.ERROR_BAD_REQUEST)
		return
	}
	if err := model.DeleteComment(common.Context, common.Client)(commentID); err != nil {
		log.Warn("failing deleting comment", zap.Int("commentID", commentID), zap.Error(err))
		c.Error(e.ERROR_BAD_REQUEST)
		return
	}
	c.Status(http.StatusOK)
}

func GetArticleComments(c *gin.Context) {
	articleID, err := utils.GetID(c, "articleID")
	if err != nil {
		return
	}

	p := page.NewPagination()
	if err := c.BindQuery(p); err != nil {
		return
	}

	comments, p, err := service.GetPublishedArticleComments(articleID, p)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments":   comments,
		"pagination": p,
	})
}
