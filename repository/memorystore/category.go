package memorystore

import (
	"github.com/SoroushBeigi/todo-cli-go/entity"
)

type Category struct {
	categories []entity.Category
}

func (c *Category) CreateNewCategory(cat entity.Category) (entity.Category, error) {
	cat.ID = len(c.categories) + 1

	c.categories = append(c.categories, cat)

	return cat, nil
}

func (c *Category) CanCreateTaskInCategory(userID, categoryID int) bool {
	isFound := false
	for _, c := range c.categories {
		if c.ID == categoryID && c.UserID == userID {
			isFound = true

			break
		}
	}

	return isFound
}
