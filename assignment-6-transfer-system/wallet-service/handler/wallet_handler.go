package handler

import (
	"context"
	"wallet/entity"
	pb "wallet/wallet-service/proto" // Adjust the import path to your proto package
	"wallet/wallet-service/service"
)

type WalletHandler struct {
	service service.WalletService // This uses your custom WalletService interface
	pb.UnimplementedWalletServiceServer
}

func NewWalletHandler(service service.WalletService) pb.WalletServiceServer {
	return &WalletHandler{
		service: service,
	}
}

func (h *WalletHandler) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.WalletResponse, error) {
	wallet := &entity.Wallet{
		UserID:  req.UserId,
		Name:    req.Name,
		Balance: req.Balance,
	}
	createdWallet, err := h.service.CreateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}
	return &pb.WalletResponse{
		Id:      createdWallet.ID,
		UserId:  createdWallet.UserID,
		Name:    createdWallet.Name,
		Balance: createdWallet.Balance,
	}, nil
}

func (h *WalletHandler) GetWalletByID(ctx context.Context, req *pb.GetWalletByIDRequest) (*pb.WalletResponse, error) {
	wallet, err := h.service.GetWalletByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.WalletResponse{
		Id:      wallet.ID,
		UserId:  wallet.UserID,
		Name:    wallet.Name,
		Balance: wallet.Balance,
	}, nil
}

func (h *WalletHandler) UpdateWallet(ctx context.Context, req *pb.UpdateWalletRequest) (*pb.WalletResponse, error) {
	wallet := &entity.Wallet{
		ID:      req.Id,
		UserID:  req.UserId,
		Name:    req.Name,
		Balance: req.Balance,
	}
	updatedWallet, err := h.service.UpdateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}
	return &pb.WalletResponse{
		Id:      updatedWallet.ID,
		UserId:  updatedWallet.UserID,
		Name:    updatedWallet.Name,
		Balance: updatedWallet.Balance,
	}, nil
}

func (h *WalletHandler) DeleteWallet(ctx context.Context, req *pb.DeleteWalletRequest) (*pb.DeleteWalletResponse, error) {
	err := h.service.DeleteWallet(ctx, req.Id)
	if err != nil {
		return &pb.DeleteWalletResponse{Success: false}, err
	}
	return &pb.DeleteWalletResponse{Success: true}, nil
}

func (h *WalletHandler) TransferBetweenWallets(ctx context.Context, req *pb.TransferBetweenWalletsRequest) (*pb.TransferBetweenWalletsResponse, error) {
	err := h.service.TransferBetweenWallets(ctx, req.FromWalletId, req.ToWalletId, req.Amount)
	if err != nil {
		return &pb.TransferBetweenWalletsResponse{Success: false}, err
	}
	return &pb.TransferBetweenWalletsResponse{Success: true}, nil
}
