package model

import (
	"github.com/georgysavva/scany/pgxscan"
)

type Category struct {
	BaseModel
	Name string `json:"name" db:"name" binding:"required"`
}

func GetCategory(ID uint64) (Category, error) {
	var category Category
	err := pgxscan.Get(ctx, db, &category, "SELECT id, name, created_at, updated_at FROM category WHERE id = $1", ID)
	return category, err
}

func CreateCategory(userID uint64, category Category) (Category, error) {
	err := pgxscan.Get(
		ctx, db, &category,
		"INSERT INTO category (name) VALUES ($1) RETURNING *",
		category.Name,
	)
	return category, err
}

func GetCategories() ([]Category, error) {
	categories := []Category{}
	err := pgxscan.Select(ctx, db, &categories, "SELECT * FROM category")
	return categories, err
}

func UpdateCategory(category Category) error {
	_, err := db.Exec(ctx, "UPDATE category SET name = $1 WHERE id = $2", category.Name, category.ID)
	return err
}

func DeleteCategory(categoryID uint64) error {
	_, err := db.Exec(ctx, "DELETE FROM category WHERE id = $1", categoryID)
	return err
}

func CategoriesCount() (uint64, error) {
	var count uint64
	err := pgxscan.Get(ctx, db, &count, "SELECT COUNT(*) FROM category")
	return count, err
}
