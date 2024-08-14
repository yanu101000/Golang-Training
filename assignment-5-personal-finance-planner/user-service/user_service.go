package service

import (
	"financial/entity"
	"fmt"
)

type UserService interface {
	CreateUser(user *entity.User) (*entity.User, error)
	GetUserByID(id int64) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id int64) error
}

type userService struct {
	users  map[int64]*entity.User
	nextID int64
}

func NewUserService() UserService {
	return &userService{
		users:  make(map[int64]*entity.User),
		nextID: 1,
	}
}

func (s *userService) CreateUser(user *entity.User) (*entity.User, error) {
	user.ID = s.nextID
	s.users[s.nextID] = user
	s.nextID++
	return user, nil
}

func (s *userService) GetUserByID(id int64) (*entity.User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *userService) UpdateUser(user *entity.User) (*entity.User, error) {
	_, exists := s.users[user.ID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	s.users[user.ID] = user
	return user, nil
}

func (s *userService) DeleteUser(id int64) error {
	_, exists := s.users[id]
	if !exists {
		return fmt.Errorf("user not found")
	}
	delete(s.users, id)
	return nil
}
