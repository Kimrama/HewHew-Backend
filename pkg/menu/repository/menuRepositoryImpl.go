package repository

import (
	"hewhew-backend/database"
)
type MenuRepositoryImpl struct {
    db database.Database
}
func NewMenuRepositoryImpl(db database.Database) MenuRepository {
return &MenuRepositoryImpl{
    db: db,
}
}

func (r *MenuRepositoryImpl) CreateMenu(menuModel interface{}) error {
    // Implement the logic to save the menu item to the database
    // This is a placeholder implementation
    return nil
}
