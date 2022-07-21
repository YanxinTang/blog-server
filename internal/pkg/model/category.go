package model

import (
	"context"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/pkg/errors"
)

var CategoryNil = 0

func CreateCategory(ctx context.Context, client *ent.Client) func(name string) (*ent.Category, error) {
	return func(name string) (*ent.Category, error) {
		c, err := client.Category.
			Create().
			SetName(name).
			Save(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failing creating category")
		}
		return c, nil
	}
}

func DeleteCategory(ctx context.Context, client *ent.Client) func(id int) error {
	return func(id int) error {
		err := client.Category.DeleteOneID(id).Exec(ctx)
		if err != nil {
			return errors.Wrapf(err, "failing deleting categoy[%d]", id)
		}
		return nil
	}
}

type UpdateCategoryInput struct {
	ID   int
	Name string
}

func UpdateCategory(ctx context.Context, client *ent.Client) func(UpdateCategoryInput) (*ent.Category, error) {
	return func(updateCategoryInput UpdateCategoryInput) (*ent.Category, error) {
		update := client.Category.UpdateOneID(updateCategoryInput.ID)
		if updateCategoryInput.Name != "" {
			update.SetName(updateCategoryInput.Name)
		}
		c, err := update.Save(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failing updating category")
		}
		return c, nil
	}
}

func GetCategory(ctx context.Context, client *ent.Client) func(id int) (*ent.Category, error) {
	return func(id int) (*ent.Category, error) {
		c, err := client.Category.Get(ctx, id)
		if err != nil {
			return nil, errors.Wrapf(err, "failing getting category[%d]", id)
		}
		return c, nil
	}
}

func GetCategories(ctx context.Context, client *ent.Client) func() ([]*ent.Category, error) {
	return func() ([]*ent.Category, error) {
		categories, err := client.Category.Query().All(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failing getting all categories")
		}
		return categories, nil
	}
}

func CategoriesCount(ctx context.Context, client *ent.Client) func() (int, error) {
	return func() (int, error) {
		count, err := client.Category.Query().Count(ctx)
		if err != nil {
			return 0, errors.Wrap(err, "failing getting count of category")
		}
		return count, nil
	}
}
