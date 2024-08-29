package repository

import (
	"wallet/entity"

	"gorm.io/gorm"
)

type TransactionCategoryRepository interface {
	CreateTransactionCategory(category *entity.TransactionCategory) (*entity.TransactionCategory, error)
	GetTransactionCategoryByID(id int64) (*entity.TransactionCategory, error)
	UpdateTransactionCategory(category *entity.TransactionCategory) (*entity.TransactionCategory, error)
	DeleteTransactionCategory(id int64) error
}

type transactionCategoryRepository struct {
	db *gorm.DB
}

func NewTransactionCategoryRepository(db *gorm.DB) TransactionCategoryRepository {
	return &transactionCategoryRepository{db}
}

func (r *transactionCategoryRepository) CreateTransactionCategory(category *entity.TransactionCategory) (*entity.TransactionCategory, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *transactionCategoryRepository) GetTransactionCategoryByID(id int64) (*entity.TransactionCategory, error) {
	var category entity.TransactionCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *transactionCategoryRepository) UpdateTransactionCategory(category *entity.TransactionCategory) (*entity.TransactionCategory, error) {
	if err := r.db.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *transactionCategoryRepository) DeleteTransactionCategory(id int64) error {
	if err := r.db.Delete(&entity.TransactionCategory{}, id).Error; err != nil {
		return err
	}
	return nil
}
