package model

import (
	"log"
)

type Category struct {
	BaseModel
	Name string `json:"name" db:"name" binding:"required"`
}

func GetCategory(ID uint64) (Category, error) {
	var category Category
	row := DB.QueryRowx("SELECT `id`, `name`, `created_at`, `updated_at` FROM category WHERE id = ?", ID)
	if err := row.StructScan(&category); err != nil {
		return category, err
	}
	return category, nil
}

func CreateCategory(userID uint64, category Category) (Category, error) {
	res, err := DB.Exec(
		"INSERT category SET name= ?",
		category.Name,
	)
	if err != nil {
		return category, err
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return category, err
	}
	lastInsertCategory, err := GetCategory(uint64(lastInsertID))
	if err != nil {
		return category, err
	}
	return lastInsertCategory, nil
}

func GetCategories() ([]Category, error) {
	rows, err := DB.Queryx("SELECT * FROM category")
	if err != nil {
		return nil, err
	}

	categories := []Category{}
	for rows.Next() {
		var category Category
		err := rows.StructScan(&category)
		if err != nil {
			log.Println("获取全部分类出错：", err)
			continue
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func UpdateCategory(category Category) error {
	_, err := DB.Exec("UPDATE category SET name = ? WHERE id = ?", category.Name, category.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCategory(categoryID uint64) error {
	if _, err := DB.Exec("DELETE FROM category WHERE id = ?", categoryID); err != nil {
		return err
	}
	return nil
}

func CategoriesCount() uint64 {
	var count uint64
	DB.QueryRow("SELECT COUNT(*) FROM category").Scan(&count)
	return count
}
