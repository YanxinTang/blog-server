package model

import (
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/utils"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

const (
	StatusAll       = iota - 1
	StatusPublished = iota - 1
	StatusDraft     = iota - 1
)

type ArticleStatus int

type Article struct {
	BaseModel
	CategoryID uint64        `json:"categoryID" db:"category_id" binding:"required"`
	Title      string        `json:"title" db:"title" binding:"required"`
	Content    string        `json:"content" db:"content" binding:"required"`
	Status     ArticleStatus `json:"status" db:"status"`

	Category Category `json:"category" binding:"-"`
}

func (a *Article) Summary() string {
	return utils.Summary(a.Content)
}

func GetArticles(status ArticleStatus, pagination Pagination) ([]Article, error) {
	start := (pagination.Page - 1) * pagination.PerPage
	var rows pgx.Rows
	var err error

	switch status {
	case StatusAll:
		rows, err = db.Query(
			ctx,
			`SELECT a.id, a.category_id, a.title, a.created_at, a.updated_at, c.id, c.name, c.created_at, c.updated_at 
			FROM article as a
			LEFT JOIN category as c
			ON a.category_id = c.id
			ORDER BY a.id DESC LIMIT $1 OFFSET $2
			`,
			pagination.PerPage,
			start,
		)
	default:
		rows, err = db.Query(
			ctx,
			`SELECT a.id, a.category_id, a.title, a.created_at, a.updated_at, c.id, c.name, c.created_at, c.updated_at 
			FROM article as a
			LEFT JOIN category as c
			ON a.category_id = c.id
			WHERE a.status = $1
			ORDER BY a.id DESC LIMIT $2 OFFSET $3
			`,
			status,
			pagination.PerPage,
			start,
		)
	}

	if err != nil {
		return nil, err
	}

	var articles []Article = make([]Article, 0)
	for rows.Next() {
		var article Article
		if err := rows.Scan(
			&article.BaseModel.ID,
			&article.CategoryID,
			&article.Title,
			&article.BaseModel.CreatedAt,
			&article.BaseModel.UpdatedAt,
			&article.Category.ID,
			&article.Category.Name,
			&article.Category.CreatedAt,
			&article.Category.UpdatedAt,
		); err != nil {
			log.Warn("failed to scan article", zap.Error(err))
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// GetCategoryArticles 获取某个分类下某个状态的文章
func GetCategoryArticles(categoryID uint64, status ArticleStatus, pagination Pagination) ([]Article, error) {
	var articles []Article
	start := (pagination.Page - 1) * pagination.PerPage
	var err error

	if status == StatusAll {
		pgxscan.Select(
			ctx,
			db,
			&articles,
			`SELECT id, category_id, title, content, created_at, updated_at
			FROM article
			WHERE category_id = $1
			ORDER BY id DESC LIMIT $2 OFFSET $3
			`,
			categoryID,
			pagination.PerPage,
			start,
		)
	} else {
		pgxscan.Select(
			ctx,
			db,
			&articles,
			`SELECT id, category_id, title, content, created_at, updated_at
			FROM article
			WHERE category_id = $1 AND status = $2
			ORDER BY id DESC LIMIT $3 OFFSET $4
			`,
			categoryID,
			status,
			pagination.PerPage,
			start,
		)
	}

	if err != nil {
		return nil, err
	}

	category, err := GetCategory(categoryID)
	if err != nil {
		return nil, err
	}
	for i := range articles {
		articles[i].Category = category
	}
	return articles, nil
}

func getStatusArticle(articleID uint64, status ArticleStatus) (Article, error) {
	row := db.QueryRow(
		ctx,
		`SELECT a.id, a.category_id, a.title, a.content, a.created_at, a.updated_at, c.id, c.name, c.created_at, c.updated_at 
		FROM article as a
		LEFT JOIN category as c 
		ON a.category_id = c.id 
		WHERE a.id = $1 AND a.status = $2`,
		articleID,
		status,
	)

	var article Article
	err := row.Scan(
		&article.BaseModel.ID,
		&article.CategoryID,
		&article.Title,
		&article.Content,
		&article.BaseModel.CreatedAt,
		&article.BaseModel.UpdatedAt,
		&article.Category.ID,
		&article.Category.Name,
		&article.Category.CreatedAt,
		&article.Category.UpdatedAt,
	)

	return article, err
}

func GetPublishedArticle(articleID uint64) (Article, error) {
	return getStatusArticle(articleID, StatusPublished)
}

func GetDraft(draftID uint64) (Article, error) {
	return getStatusArticle(draftID, StatusDraft)
}

func CreateArticle(userID uint64, article Article) (Article, error) {
	err := pgxscan.Get(
		ctx,
		db,
		&article,
		"INSERT INTO article (title, category_id, content, status) VALUES ($1, $2, $3, $4) RETURNING *",
		&article.Title,
		&article.CategoryID,
		&article.Content,
		&article.Status,
	)
	return article, err
}

func DeleteArticle(articleID uint64) error {
	_, err := db.Exec(ctx, "DELETE FROM article WHERE id = $1", articleID)
	return err
}

func UpdateArticle(article Article) (pgconn.CommandTag, error) {
	return db.Exec(
		ctx,
		"UPDATE article SET category_id = $1, title = $2, content = $3 WHERE id = $4",
		article.CategoryID,
		article.Title,
		article.Content,
		article.ID,
	)
}

func GetArticlesCount(categoryID uint64, status ArticleStatus) (uint64, error) {
	var count uint64
	var err error
	switch {
	case categoryID != 0 && status != StatusAll:
		err = pgxscan.Get(
			ctx,
			db,
			&count,
			"SELECT COUNT(*) FROM article WHERE category_id = $1 AND status = $2",
			categoryID,
			status,
		)
	case categoryID != 0:
		err = pgxscan.Get(
			ctx,
			db,
			&count,
			"SELECT COUNT(*) FROM article WHERE category_id = $1",
			categoryID,
		)
	case status != StatusAll:
		err = pgxscan.Get(
			ctx,
			db,
			&count,
			"SELECT COUNT(*) FROM article WHERE status = $1",
			status,
		)
	default:
		err = pgxscan.Get(
			ctx,
			db,
			&count,
			"SELECT COUNT(*) FROM article",
		)
	}
	return count, err
}

// CategoryArticlesCount 返回某个分类下文章的总数
func CategoryArticlesCount(categoryID uint64) (uint64, error) {
	var count uint64
	err := pgxscan.Get(
		ctx,
		db,
		&count,
		"SELECT COUNT(*) FROM article WHERE status = $1 AND category_id = $2",
		StatusPublished,
		categoryID,
	)
	return count, err
}

func PublishDraft(draftID uint64) error {
	_, err := db.Exec(ctx, "UPDATE article SET status = $1 WHERE id = $2", StatusPublished, draftID)
	return err
}

func DraftsCount() (uint64, error) {
	var count uint64
	err := pgxscan.Get(ctx, db, &count, "SELECT COUNT(*) FROM article WHERE status = $1", StatusDraft)
	return count, err
}
