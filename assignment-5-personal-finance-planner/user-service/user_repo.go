package repository

import (
	"financial/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	GetUserByID(id int64) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id int64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *entity.User) (*entity.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(id int64) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) DeleteUser(id int64) error {
	if err := r.db.Delete(&entity.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
