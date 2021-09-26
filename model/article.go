package model

import (
	"database/sql"
	"log"

	"github.com/YanxinTang/blog/server/utils"
)

const (
	StatusPublished = iota
	StatusDraft     = iota
)

type Article struct {
	BaseModel
	CategoryID uint64 `json:"categoryID" db:"category_id" binding:"required"`
	Title      string `json:"title" db:"title" binding:"required"`
	Content    string `json:"content" db:"content" binding:"required"`
	Status     int8   `json:"status" db:"status"`

	Category Category `json:"category" binding:"-"`
}

func (a *Article) Summary() string {
	return utils.Summary(a.Content)
}

func getArticles(status int8, pagination Pagination) ([]Article, error) {
	start := (pagination.Page - 1) * pagination.PerPage
	rows, err := DB.Query(
		`SELECT a.id, a.category_id, a.title, a.content, a.created_at, a.updated_at, c.id, c.name, c.created_at, c.updated_at 
		FROM article as a
		LEFT JOIN category as c
		ON a.category_id = c.id
		WHERE a.status = ?
		ORDER BY a.id DESC LIMIT ?, ?
		`,
		status,
		start,
		pagination.PerPage,
	)
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
			&article.Content,
			&article.BaseModel.CreatedAt,
			&article.BaseModel.UpdatedAt,
			&article.Category.ID,
			&article.Category.Name,
			&article.Category.CreatedAt,
			&article.Category.UpdatedAt,
		); err != nil {
			log.Println("获取 articles 错误：", err)
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func GetPublishedArticles(pagination Pagination) ([]Article, error) {
	return getArticles(StatusPublished, pagination)
}

func GetDrafts(pagination Pagination) ([]Article, error) {
	return getArticles(StatusDraft, pagination)
}

// getCategoryArticles 获取某个分类下的所有内容
func getCategoryArticles(categoryID uint64, status int8, pagination Pagination) ([]Article, error) {
	var articles []Article
	start := (pagination.Page - 1) * pagination.PerPage
	err := DB.Select(
		&articles,
		`SELECT id, category_id, title, content, created_at, updated_at
		FROM article as a
		WHERE category_id = ? AND status = ?
		ORDER BY id DESC LIMIT ?, ?
		`,
		categoryID,
		status,
		start,
		pagination.PerPage,
	)
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

// 获取某个分类下的所有文章
func GetCategoryPublishedArticles(categoryID uint64, pagination Pagination) ([]Article, error) {
	return getCategoryArticles(categoryID, StatusPublished, pagination)
}

func getArticle(articleID uint64) (Article, error) {
	row := DB.QueryRow(
		`SELECT a.id, a.category_id, a.title, a.content, a.created_at, a.updated_at, c.id, c.name, c.created_at, c.updated_at 
		FROM article as a
		LEFT JOIN category as c 
		ON a.category_id = c.id 
		WHERE a.id = ?`,
		articleID,
	)

	var article Article
	if err := row.Scan(
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
	); err != nil {
		return Article{}, err
	}

	return article, nil
}

func getStatusArticle(articleID uint64, status int8) (Article, error) {
	row := DB.QueryRow(
		`SELECT a.id, a.category_id, a.title, a.content, a.created_at, a.updated_at, c.id, c.name, c.created_at, c.updated_at 
		FROM article as a
		LEFT JOIN category as c 
		ON a.category_id = c.id 
		WHERE a.id = ? AND a.status = ?`,
		articleID,
		status,
	)

	var article Article
	if err := row.Scan(
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
	); err != nil {
		return Article{}, err
	}

	return article, nil
}

func GetPublishedArticle(articleID uint64) (Article, error) {
	return getStatusArticle(articleID, StatusPublished)
}

func GetDraft(draftID uint64) (Article, error) {
	return getStatusArticle(draftID, StatusDraft)
}

func CreateArticle(userID uint64, article Article) (Article, error) {
	res, err := DB.Exec(
		"INSERT article SET title=?, category_id = ?, content=?, status = ?",
		&article.Title,
		&article.CategoryID,
		&article.Content,
		&article.Status,
	)
	if err != nil {
		return article, err
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return article, err
	}
	lastInsertArticle, err := getArticle(uint64(lastInsertID))
	if err != nil {
		return article, err
	}
	return lastInsertArticle, nil
}

func DeleteArticle(articleID uint64) error {
	if _, err := DB.Exec("DELETE FROM article WHERE id = ?", articleID); err != nil {
		return err
	}
	return nil
}

func UpdateArticle(article Article) (sql.Result, error) {
	stmt, err := DB.Prepare("UPDATE article SET category_id = ?, title = ?, content = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(article.CategoryID, article.Title, article.Content, article.ID)
}

func ArticlesCount() uint64 {
	var count uint64
	DB.QueryRow("SELECT COUNT(*) FROM article WHERE status = ?", StatusPublished).Scan(&count)
	return count
}

// CategoryArticlesCount 返回某个分类下文章的总数
func CategoryArticlesCount(categoryID uint64) uint64 {
	var count uint64
	DB.QueryRow(
		"SELECT COUNT(*) FROM article WHERE status = ? AND category_id = ?",
		StatusPublished,
		categoryID,
	).Scan(&count)
	return count
}

func PublishDraft(draftID uint64) error {
	if _, err := DB.Exec("UPDATE `article` SET `status` = ? WHERE `id` = ?", StatusPublished, draftID); err != nil {
		return err
	}
	return nil
}

func DraftsCount() uint64 {
	var count uint64
	DB.QueryRow("SELECT COUNT(*) FROM article WHERE status = ?", StatusDraft).Scan(&count)
	return count
}
