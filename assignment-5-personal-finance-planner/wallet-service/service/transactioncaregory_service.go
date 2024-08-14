package service

import (
	"fmt"
	"wallet/entity"
)

type TransactionCategoryService interface {
	CreateTransactionCategory(category *entity.TransactionCategory) (*entity.TransactionCategory, error)
	GetTransactionCategoryByID(id int64) (*entity.TransactionCategory, error)
	UpdateTransactionCategory(category *entity.TransactionCategory) (*entity.TransactionCategory, error)
	DeleteTransactionCategory(id int64) error
	GetCategories() map[int64]*entity.TransactionCategory
}

type transactionCategoryService struct {
	categories map[int64]*entity.TransactionCategory
	nextID     int64
}

func NewTransactionCategoryService() TransactionCategoryService {
	return &transactionCategoryService{
		categories: make(map[int64]*entity.TransactionCategory),
		nextID:     1,
	}
}

func (s *transactionCategoryService) CreateTransactionCategory(category *entity.TransactionCategory) (*entity.TransactionCategory, error) {
	category.ID = s.nextID
	s.categories[s.nextID] = category
	s.nextID++
	return category, nil
}

func (s *transactionCategoryService) GetTransactionCategoryByID(id int64) (*entity.TransactionCategory, error) {
	category, exists := s.categories[id]
	if !exists {
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}

func (s *transactionCategoryService) UpdateTransactionCategory(category *entity.TransactionCategory) (*entity.TransactionCategory, error) {
	_, exists := s.categories[category.ID]
	if !exists {
		return nil, fmt.Errorf("category not found")
	}
	s.categories[category.ID] = category
	return category, nil
}

func (s *transactionCategoryService) DeleteTransactionCategory(id int64) error {
	_, exists := s.categories[id]
	if !exists {
		return fmt.Errorf("category not found")
	}
	delete(s.categories, id)
	return nil
}

func (s *transactionCategoryService) GetCategories() map[int64]*entity.TransactionCategory {
	return s.categories
}
