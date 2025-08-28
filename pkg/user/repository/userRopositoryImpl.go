package repository

import (
	"hewhew-backend/database"
	"hewhew-backend/entities"
)

type UserRepositoryImpl struct {
	db database.Database
}

func NewUserRepositoryImpl(db database.Database) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) CreateUser(userEntity *entities.User) (*entities.User, error) {
	return userEntity, nil
}
func (r *UserRepositoryImpl) GetUsers() ([]*entities.User, error) {
	var users []*entities.User
	r.db.Connect()
	return users, nil
}
