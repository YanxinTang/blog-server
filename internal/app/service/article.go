package service

import "github.com/YanxinTang/blog-server/internal/pkg/model"

func GetArticles(categoryID uint64, status model.ArticleStatus, pagination model.Pagination) ([]model.Article, error) {
	if categoryID == 0 {
		return model.GetArticles(status, pagination)
	}
	return model.GetCategoryArticles(categoryID, status, pagination)
}

func GetPublishedArticles(pagination model.Pagination) ([]model.Article, error) {
	return model.GetArticles(model.StatusPublished, pagination)
}

func GetDrafts(pagination model.Pagination) ([]model.Article, error) {
	return model.GetArticles(model.StatusDraft, pagination)
}

// 获取某个分类下的已发布文章
func GetCategoryPublishedArticles(categoryID uint64, pagination model.Pagination) ([]model.Article, error) {
	return model.GetCategoryArticles(categoryID, model.StatusPublished, pagination)
}
