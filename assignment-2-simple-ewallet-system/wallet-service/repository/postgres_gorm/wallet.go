package postgres_gorm

import (
	"context"
	"errors"
	"log"
	"solution1/assignment-2-simple-ewallet-system/wallet-service/entity"
	"solution1/assignment-2-simple-ewallet-system/wallet-service/service"

	"gorm.io/gorm"
)

// GormDBIface defines an interface for GORM DB methods used in the repository
type GormDBIface interface {
	WithContext(ctx context.Context) *gorm.DB
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
}
type walletRepository struct {
	db GormDBIface
}

// NewWalletRepository creates a new instance of walletRepository
func NewWalletRepository(db GormDBIface) service.IWalletRepository {
	return &walletRepository{db: db}
}

// CreateWallet creates a new wallet in the database
func (r *walletRepository) CreateWallet(ctx context.Context, userID int32) (*entity.Wallet, error) {
	wallet := &entity.Wallet{
		UserID:  uint(userID), // Convert int32 to uint, if UserID in entity.Wallet is uint
		Balance: 0.0,          // Initialize balance as needed
	}

	if err := r.db.WithContext(ctx).Create(wallet).Error; err != nil {
		log.Printf("Error creating wallet: %v\n", err)
		return nil, err
	}
	return wallet, nil
}

// GetWallet retrieves a wallet by user ID
func (r *walletRepository) GetWallet(ctx context.Context, userID int32) (*entity.Wallet, error) {
	var wallet entity.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("Error getting wallet for user %d: %v\n", userID, err)
		return nil, err
	}
	return &wallet, nil
}

// TopUp adds an amount to a user's wallet
func (r *walletRepository) TopUp(ctx context.Context, userID int32, amount float32) (*entity.Wallet, error) {
	wallet, err := r.GetWallet(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	wallet.Balance += amount
	if err := r.db.WithContext(ctx).Save(wallet).Error; err != nil {
		log.Printf("Error topping up wallet for user %d: %v\n", userID, err)
		return nil, err
	}
	return wallet, nil
}

// Transfer performs a transfer of amount from one user's wallet to another
func (r *walletRepository) Transfer(ctx context.Context, fromUserID, toUserID int32, amount float32) (*entity.Wallet, error) {
	fromWallet, err := r.GetWallet(ctx, fromUserID)
	if err != nil {
		return nil, err
	}
	if fromWallet == nil {
		return nil, errors.New("from user's wallet not found")
	}

	toWallet, err := r.GetWallet(ctx, toUserID)
	if err != nil {
		return nil, err
	}
	if toWallet == nil {
		return nil, errors.New("to user's wallet not found")
	}

	if fromWallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	fromWallet.Balance -= amount
	toWallet.Balance += amount

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(fromWallet).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating from user's wallet: %v\n", err)
		return nil, err
	}

	if err := tx.Save(toWallet).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating to user's wallet: %v\n", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return nil, err
	}

	return fromWallet, nil
}

// GetTransactions retrieves transactions for a user
func (r *walletRepository) GetTransactions(ctx context.Context, userID int32) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactions, nil
		}
		log.Printf("Error getting transactions for user %d: %v\n", userID, err)
		return nil, err
	}
	return transactions, nil
}
