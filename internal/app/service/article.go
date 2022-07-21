package service

import (
	"net/http"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/YanxinTang/blog-server/internal/pkg/page"
	"github.com/YanxinTang/blog-server/utils"
	"go.uber.org/zap"
)

func DeleteArticle(articleID int) error {
	return utils.WithTx(common.Context, common.Client, func(tx *ent.Tx) error {
		if err := model.DeleteArticleComments(common.Context, common.Client)(articleID); err != nil {
			return err
		}
		if err := model.DeleteArticle(common.Context, tx.Client())(articleID); err != nil {
			return err
		}
		return nil
	})
}

func GetArticles(categoryID int, status model.ArticleStatus, pagination *page.Pagination) ([]*ent.Article, *page.Pagination, error) {
	if categoryID == 0 {
		return model.GetArticles(common.Context, common.Client)(status, pagination)
	}
	return model.GetCategoryArticles(common.Context, common.Client)(categoryID, status, pagination)
}

func GetPublishedArticles(pagination *page.Pagination) ([]*ent.Article, *page.Pagination, error) {
	return model.GetArticles(common.Context, common.Client)(model.StatusPublished, pagination)
}

func GetDrafts(p *page.Pagination) ([]*ent.Article, *page.Pagination, error) {
	return model.GetArticles(common.Context, common.Client)(model.StatusPublished, p)
}

// 获取某个分类下的已发布文章
func GetCategoryPublishedArticles(categoryID int, p *page.Pagination) ([]*ent.Article, *page.Pagination, error) {
	return model.GetCategoryArticles(common.Context, common.Client)(categoryID, model.StatusPublished, p)
}

func GetPublishedArticle(articleID int) (*ent.Article, error) {
	article, err := model.GetArticle(common.Context, common.Client)(articleID, model.StatusPublished)
	if err != nil {
		log.Warn("failing getting published article", zap.Int("articleID", articleID), zap.Error(err))
		return nil, e.New(http.StatusNotFound, "文章未找到")
	}
	return article, nil
}
