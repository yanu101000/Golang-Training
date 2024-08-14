package service

import (
	"context"
	"fmt"
	"wallet/entity"
	"wallet/repository"
)

type WalletService interface {
	CreateWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error)
	GetWalletByID(ctx context.Context, id int64) (*entity.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error)
	DeleteWallet(ctx context.Context, id int64) error
	TransferBetweenWallets(ctx context.Context, fromWalletID, toWalletID int64, amount float64) error
}

type walletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{
		repo: repo,
	}
}

func (s *walletService) CreateWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
	createdWallet, err := s.repo.CreateWallet(wallet)
	if err != nil {
		return nil, fmt.Errorf("could not create wallet: %w", err)
	}
	return createdWallet, nil
}

func (s *walletService) GetWalletByID(ctx context.Context, id int64) (*entity.Wallet, error) {
	wallet, err := s.repo.GetWalletByID(id)
	if err != nil {
		return nil, fmt.Errorf("could not get wallet: %w", err)
	}
	return wallet, nil
}

func (s *walletService) UpdateWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
	updatedWallet, err := s.repo.UpdateWallet(wallet)
	if err != nil {
		return nil, fmt.Errorf("could not update wallet: %w", err)
	}
	return updatedWallet, nil
}

func (s *walletService) DeleteWallet(ctx context.Context, id int64) error {
	if err := s.repo.DeleteWallet(id); err != nil {
		return fmt.Errorf("could not delete wallet: %w", err)
	}
	return nil
}

func (s *walletService) TransferBetweenWallets(ctx context.Context, fromWalletID, toWalletID int64, amount float64) error {
	if err := s.repo.TransferBetweenWallets(fromWalletID, toWalletID, amount); err != nil {
		return fmt.Errorf("could not transfer between wallets: %w", err)
	}
	return nil
}
