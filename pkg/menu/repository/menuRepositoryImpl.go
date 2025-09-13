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


