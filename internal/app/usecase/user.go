package usecase

import (
	"file_storage_service/internal/domain/model"
	"file_storage_service/internal/domain/repository"
)

type UserUseCase struct {
	UserRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepo: userRepo,
	}
}

func (uc *UserUseCase) GetUsers() ([]*model.User, error) {
	return uc.UserRepo.GetUsers()
}
