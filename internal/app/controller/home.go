package controller

import (
	"net/http"

	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/app/service"
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/gin-gonic/gin"
)

type Pair struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func Overview(c *gin.Context) {
	articlesCount, err := model.ArticlesCount(common.Context, common.Client)()
	if err != nil {
		c.Error(e.New(http.StatusInternalServerError, "获取文章数量失败"))
		return
	}

	categoriesCount, err := model.CategoriesCount(common.Context, common.Client)()
	if err != nil {
		c.Error(e.New(http.StatusInternalServerError, "获取分类数量失败"))
		return
	}

	commentsCount, err := model.CommentsCount(common.Context, common.Client)()
	if err != nil {
		c.Error(e.New(http.StatusInternalServerError, "获取评论数量失败"))
		return
	}

	pairs := make([]Pair, 0, 3)
	pairs = append(
		pairs,
		Pair{"分类数量", categoriesCount},
		Pair{"文章数量", articlesCount},
		Pair{"评论数量", commentsCount},
	)
	c.JSON(http.StatusOK, pairs)
}

type StorageOverviewItem struct {
	ID       int    `json:"id"`
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
