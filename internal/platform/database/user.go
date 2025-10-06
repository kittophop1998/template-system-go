package database

import (
	"file_storage_service/internal/domain/model"

	"gorm.io/gorm"
)

type UserPostGres struct {
	db *gorm.DB
}

func NewUserPostGres(db *gorm.DB) *UserPostGres {
	return &UserPostGres{db: db}
}

func (r *UserPostGres) GetUsers() ([]*model.User, error) {
	var users []*model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
