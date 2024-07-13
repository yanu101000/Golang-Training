package service

import (
	"context"
	"errors"
	"log"
	"solution1/assignment-2-simple-ewallet-system/wallet-service/entity"
)

// WalletService defines the methods for handling wallet operations.
type IWalletService interface {
	CreateWallet(ctx context.Context, userID int32) (*entity.Wallet, error)
	TopUp(ctx context.Context, userID int32, amount float32) (*entity.Wallet, error)
	Transfer(ctx context.Context, fromUserID, toUserID int32, amount float32) (*entity.Wallet, error)
	GetWallet(ctx context.Context, userID int32) (*entity.Wallet, error)
	GetTransactions(ctx context.Context, userID int32) ([]*entity.Transaction, error)
}

// WalletRepository defines the methods for data access related to wallets.
type IWalletRepository interface {
	CreateWallet(ctx context.Context, userID int32) (*entity.Wallet, error)
	GetWallet(ctx context.Context, userID int32) (*entity.Wallet, error)
	TopUp(ctx context.Context, userID int32, amount float32) (*entity.Wallet, error)
	Transfer(ctx context.Context, fromUserID, toUserID int32, amount float32) (*entity.Wallet, error)
	GetTransactions(ctx context.Context, userID int32) ([]*entity.Transaction, error)
}

type WalletServiceImpl struct {
	repo IWalletRepository
}

func NewWalletService(repo IWalletRepository) IWalletService {
	return &WalletServiceImpl{repo: repo}
}

func (s *WalletServiceImpl) CreateWallet(ctx context.Context, userID int32) (*entity.Wallet, error) {
	// Implement logic to create a new wallet
	wallet, err := s.repo.CreateWallet(ctx, userID)
	if err != nil {
		log.Printf("failed to create wallet for user %d: %v", userID, err)
		return nil, errors.New("failed to create wallet")
	}
	return wallet, nil
}

func (s *WalletServiceImpl) TopUp(ctx context.Context, userID int32, amount float32) (*entity.Wallet, error) {
	// Implement logic to top up a wallet
	wallet, err := s.repo.TopUp(ctx, userID, amount)
	if err != nil {
		log.Printf("failed to top up wallet for user %d: %v", userID, err)
		return nil, errors.New("failed to top up wallet")
	}
	return wallet, nil
}

func (s *WalletServiceImpl) Transfer(ctx context.Context, fromUserID, toUserID int32, amount float32) (*entity.Wallet, error) {
	// Implement logic to transfer amount between wallets
	fromWallet, err := s.repo.GetWallet(ctx, fromUserID)
	if err != nil {
		log.Printf("failed to retrieve wallet for user %d: %v", fromUserID, err)
		return nil, errors.New("failed to retrieve wallet")
	}

	if fromWallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	toWallet, err := s.repo.GetWallet(ctx, toUserID)
	if err != nil {
		log.Printf("failed to retrieve wallet for user %d: %v", toUserID, err)
		return nil, errors.New("failed to retrieve wallet")
	}

	fromWallet.Balance -= amount
	toWallet.Balance += amount

	// Update balances in the repository
	_, err = s.repo.Transfer(ctx, fromUserID, toUserID, amount)
	if err != nil {
		log.Printf("failed to transfer amount from user %d to user %d: %v", fromUserID, toUserID, err)
		return nil, errors.New("failed to transfer amount")
	}

	return fromWallet, nil
}

func (s *WalletServiceImpl) GetWallet(ctx context.Context, userID int32) (*entity.Wallet, error) {
	// Implement logic to retrieve a wallet
	wallet, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		log.Printf("failed to retrieve wallet for user %d: %v", userID, err)
		return nil, errors.New("failed to retrieve wallet")
	}
	return wallet, nil
}

func (s *WalletServiceImpl) GetTransactions(ctx context.Context, userID int32) ([]*entity.Transaction, error) {
	// Implement logic to retrieve transactions for a user
	transactions, err := s.repo.GetTransactions(ctx, userID)
	if err != nil {
		log.Printf("failed to retrieve transactions for user %d: %v", userID, err)
		return nil, errors.New("failed to retrieve transactions")
	}
	return transactions, nil
}
