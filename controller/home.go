package controller

import (
	"net/http"
	"strconv"

	"github.com/YanxinTang/blog-server/model"
	"github.com/YanxinTang/blog-server/service"
	"github.com/gin-gonic/gin"
)

type Pair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func Overview(c *gin.Context) {
	articlesCount, err := model.ArticlesCount()
	if err != nil {
		c.Error(err)
		return
	}

	categoriesCount, err := model.CategoriesCount()
	if err != nil {
		c.Error(err)
		return
	}

	commentsCount, err := model.CommentsCount()
	if err != nil {
		c.Error(err)
		return
	}

	pairs := make([]Pair, 0, 3)
	pairs = append(
		pairs,
		Pair{"分类数量", strconv.FormatUint(categoriesCount, 10)},
		Pair{"文章数量", strconv.FormatUint(articlesCount, 10)},
		Pair{"评论数量", strconv.FormatUint(commentsCount, 10)},
	)
	c.JSON(http.StatusOK, pairs)
}

type StorageOverviewItem struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Usage    int64  `json:"usage"`
	Capacity int64  `json:"capacity"`
}

func OverviewStorage(c *gin.Context) {
	svcs, err := service.GetStorageServices()
	if err != nil {
		c.Error(err)
		return
	}

	storageOverviewItems := make([]StorageOverviewItem, 0, len(svcs))
	for _, svc := range svcs {
		usage, err := service.GetStorageUsage(svc)
		if err != nil {
			c.Error(err)
			return
		}
		storageOverviewItems = append(storageOverviewItems, StorageOverviewItem{
			ID:       svc.Storage.ID,
			Name:     svc.Storage.Name,
			Usage:    usage,
			Capacity: svc.Storage.Capacity,
		})
	}
	c.JSON(http.StatusOK, storageOverviewItems)
}
