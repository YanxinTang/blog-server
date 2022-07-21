package model

import (
	"context"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/ent/article"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/internal/pkg/page"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ArticleStatus int8

const (
	StatusNil = ArticleStatus(iota)
	StatusPublished
	StatusDraft
)

type CreateArticleInput struct {
	CategoryID int
	Title      string
	Content    string
	Status     ArticleStatus
}

func CreateArticle(ctx context.Context, client *ent.Client) func(CreateArticleInput) (*ent.Article, error) {
	return func(createArticleInput CreateArticleInput) (*ent.Article, error) {
		article, err := client.Article.
			Create().
			SetCategoryID(createArticleInput.CategoryID).
			SetTitle(createArticleInput.Title).
			SetContent(createArticleInput.Content).
			SetStatus(int8(createArticleInput.Status)).
			Save(ctx)
		if err != nil {
			log.Warn("failing create article", zap.Error(err))
			return nil, errors.Wrapf(err, "failing create article")
		}
		return article, nil
	}
}

func DeleteArticle(ctx context.Context, client *ent.Client) func(articleID int) error {
	return func(articleID int) error {
		err := client.Article.DeleteOneID(articleID).Exec(ctx)
		if err != nil {
			log.Warn("failing delete article", zap.Int("articleID", articleID), zap.Error(err))
			return errors.Wrapf(err, "failing deleting article[%d]", articleID)
		}
		return nil
	}
}

type UpdateArticleInput struct {
	ID         int
	CategoryID int
	Title      string
	Content    string
	Status     ArticleStatus
}

func UpdateArticle(ctx context.Context, client *ent.Client) func(UpdateArticleInput) (*ent.Article, error) {
	return func(updateArticleInput UpdateArticleInput) (*ent.Article, error) {
		articleClient := client.Article.UpdateOneID(updateArticleInput.ID)
		if updateArticleInput.Title != "" {
			articleClient.SetTitle(updateArticleInput.Title)
		}
		if updateArticleInput.Content != "" {
			articleClient.SetContent(updateArticleInput.Content)
		}
		if updateArticleInput.CategoryID > 0 {
			articleClient.SetCategoryID(updateArticleInput.CategoryID)
		}
		if updateArticleInput.Status > StatusNil {
			articleClient.SetStatus(int8(updateArticleInput.Status))
		}
		a, err := articleClient.Save(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failing updating article")
		}
		return a, nil
	}
}

func GetArticle(ctx context.Context, client *ent.Client) func(id int, status ArticleStatus) (*ent.Article, error) {
	return func(id int, status ArticleStatus) (*ent.Article, error) {
		query := client.Article.Query().Where(article.ID(id))
		if status != StatusNil {
			query = query.Where(article.Status(int8(status)))
		}
		a, err := query.WithCategory().Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "failing getting the article[%d] with status [%d]", id, status)
		}
		return a, nil
	}
}

func GetArticles(ctx context.Context, client *ent.Client) func(ArticleStatus, *page.Pagination) ([]*ent.Article, *page.Pagination, error) {
	return func(status ArticleStatus, p *page.Pagination) ([]*ent.Article, *page.Pagination, error) {
		query := client.Article.Query()
		if status != StatusNil {
			query.Where(article.Status(int8(status)))
		}
		articles, err := query.
			Clone().
			WithCategory().
			Offset(p.Offset()).Limit(p.Limit()).
			All(ctx)
		if err != nil {
			log.Warn("failing getting articles", zap.Int8("status", int8(status)), zap.Error(err))
			return nil, nil, err
		}
		count, err := query.Clone().Count(ctx)
		if err != nil {
			log.Warn("failing getting count of articles", zap.Int8("status", int8(status)), zap.Error(err))
			return nil, nil, err
		}
		p.Total = count
		return articles, p, nil
	}
}

// GetCategoryArticles 获取某个分类下某个状态的文章
func GetCategoryArticles(ctx context.Context, client *ent.Client) func(categoryID int, status ArticleStatus, p *page.Pagination) ([]*ent.Article, *page.Pagination, error) {
	return func(categoryID int, status ArticleStatus, p *page.Pagination) ([]*ent.Article, *page.Pagination, error) {
		category, err := GetCategory(ctx, client)(categoryID)
		if err != nil {
			return nil, nil, err
		}
		articles, err := category.QueryArticles().WithCategory().Offset(p.Offset()).Limit(p.Limit()).All(ctx)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failing getting articles in category[%d]", categoryID)
		}

		count, err := GetCategoryArticlesCount(ctx, client)(categoryID, status)
		if err != nil {
			return nil, nil, err
		}
		p.Total = count

		return articles, p, nil
	}
}

func GetCategoryArticlesCount(ctx context.Context, client *ent.Client) func(categoryID int, status ArticleStatus) (int, error) {
	return func(categoryID int, status ArticleStatus) (int, error) {
		category, err := GetCategory(ctx, client)(categoryID)
		if err != nil {
			return 0, err
		}
		query := category.QueryArticles()
		if status != StatusNil {
			query = query.Where(article.Status(int8(status)))
		}
		count, err := query.Count(ctx)
		if err != nil {
			log.Warn("failing getting count of articles in category", zap.Int("categoryID", categoryID), zap.Int8("status", int8(status)))
			return 0, err
		}
		log.Debug("category articles count", zap.Int("count", count))
		return count, nil
	}
}

func ArticlesCount(ctx context.Context, client *ent.Client) func() (int, error) {
	return func() (int, error) {
		count, err := client.Article.Query().Count(ctx)
		if err != nil {
			return 0, errors.Wrap(err, "failing getting count of article")
		}
		return count, nil
	}
}
