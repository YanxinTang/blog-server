package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/YanxinTang/blog/server/model"
	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("articleID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	var comment model.Comment
	comment.ArticleID = articleID
	if err := c.Bind(&comment); err != nil {
		return
	}

	comment.Username = strings.TrimSpace(comment.Username)
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
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}
