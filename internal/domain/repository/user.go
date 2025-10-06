package repository

import "file_storage_service/internal/domain/model"

type UserRepository interface {
	GetUsers() ([]*model.User, error)
}
