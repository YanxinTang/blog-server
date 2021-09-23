package controller

import (
	"net/http"
	"strconv"

	"github.com/YanxinTang/blog/server/model"
	"github.com/gin-gonic/gin"
)

type Pair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func Overview(c *gin.Context) {
	// session := sessions.Default(c)

	articlesCount := model.ArticlesCount()
	categoriesCount := model.CategoriesCount()
	commentsCount := model.CommentsCount()

	pairs := make([]Pair, 0, 3)
	pairs = append(
		pairs,
		Pair{"分类数量", strconv.FormatUint(categoriesCount, 10)},
		Pair{"文章数量", strconv.FormatUint(articlesCount, 10)},
		Pair{"评论数量", strconv.FormatUint(commentsCount, 10)},
	)
	c.JSON(http.StatusOK, pairs)
}
