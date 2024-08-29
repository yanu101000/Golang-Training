package repository

import (
	"errors"
	"wallet/entity"

	"gorm.io/gorm"
)

type WalletRepository interface {
	CreateWallet(wallet *entity.Wallet) (*entity.Wallet, error)
	GetWalletByID(id int64) (*entity.Wallet, error)
	UpdateWallet(wallet *entity.Wallet) (*entity.Wallet, error)
	DeleteWallet(id int64) error
	TransferBetweenWallets(fromWalletID, toWalletID int64, amount float64) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db}
}

func (r *walletRepository) CreateWallet(wallet *entity.Wallet) (*entity.Wallet, error) {
	if err := r.db.Create(wallet).Error; err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *walletRepository) GetWalletByID(id int64) (*entity.Wallet, error) {
	var wallet entity.Wallet
	if err := r.db.First(&wallet, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) UpdateWallet(wallet *entity.Wallet) (*entity.Wallet, error) {
	if err := r.db.Save(wallet).Error; err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *walletRepository) DeleteWallet(id int64) error {
	if err := r.db.Delete(&entity.Wallet{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wallet not found")
		}
		return err
	}
	return nil
}

func (r *walletRepository) TransferBetweenWallets(fromWalletID, toWalletID int64, amount float64) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var fromWallet, toWallet entity.Wallet

	// Fetch fromWallet
	if err := tx.First(&fromWallet, fromWalletID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("source wallet not found")
		}
		return err
	}

	// Fetch toWallet
	if err := tx.First(&toWallet, toWalletID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("destination wallet not found")
		}
		return err
	}

	// Check for sufficient funds
	if fromWallet.Balance < amount {
		tx.Rollback()
		return errors.New("insufficient funds")
	}

	// Perform the transfer
	fromWallet.Balance -= amount
	toWallet.Balance += amount

	// Save both wallets
	if err := tx.Save(&fromWallet).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(&toWallet).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
