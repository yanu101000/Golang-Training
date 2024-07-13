// wallet_handler.go

package grpc

import (
	"context"
	"fmt"
	"log"
	pb "solution1/assignment-2-simple-ewallet-system/wallet-service/proto/wallet_service/v1"
	"solution1/assignment-2-simple-ewallet-system/wallet-service/service"
)

// WalletHandler implements the WalletServiceServer interface
type WalletHandler struct {
	pb.UnimplementedWalletServiceServer
	walletService service.IWalletService
}

// NewWalletHandler creates a new instance of WalletHandler
func NewWalletHandler(walletService service.IWalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (h *WalletHandler) CreateWallet(ctx context.Context, req *pb.WalletRequest) (*pb.WalletResponse, error) {
	userID := req.GetUserId()

	wallet, err := h.walletService.CreateWallet(ctx, userID)
	if err != nil {
		log.Printf("Failed to create wallet for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to create wallet: %v", err)
	}

	return &pb.WalletResponse{
		Wallet: &pb.Wallet{
			Id:      int32(wallet.WalletID),
			UserId:  int32(wallet.UserID),
			Balance: wallet.Balance,
		},
	}, nil
}

func (h *WalletHandler) TopUp(ctx context.Context, req *pb.TopUpRequest) (*pb.TopUpResponse, error) {
	userID := req.GetUserId()
	amount := req.GetAmount()

	wallet, err := h.walletService.TopUp(ctx, userID, amount)
	if err != nil {
		log.Printf("Failed to top up wallet for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to top up wallet: %v", err)
	}

	return &pb.TopUpResponse{
		Wallet: &pb.Wallet{
			Id:      int32(wallet.WalletID),
			UserId:  int32(wallet.UserID),
			Balance: wallet.Balance,
		},
	}, nil
}

func (h *WalletHandler) Transfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	fromUserID := req.GetFromUserId()
	toUserID := req.GetToUserId()
	amount := req.GetAmount()

	wallet, err := h.walletService.Transfer(ctx, fromUserID, toUserID, amount)
	if err != nil {
		log.Printf("Failed to transfer amount from user %d to user %d: %v", fromUserID, toUserID, err)
		return nil, fmt.Errorf("failed to transfer amount: %v", err)
	}

	return &pb.TransferResponse{
		Wallet: &pb.Wallet{
			Id:      int32(wallet.WalletID),
			UserId:  int32(wallet.UserID),
			Balance: wallet.Balance,
		},
	}, nil
}

func (h *WalletHandler) GetWallet(ctx context.Context, req *pb.GetWalletRequest) (*pb.GetWalletResponse, error) {
	userID := req.GetUserId()

	wallet, err := h.walletService.GetWallet(ctx, userID)
	if err != nil {
		log.Printf("Failed to retrieve wallet for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to retrieve wallet: %v", err)
	}

	return &pb.GetWalletResponse{
		Wallet: &pb.Wallet{
			Id:      int32(wallet.WalletID),
			UserId:  int32(wallet.UserID),
			Balance: wallet.Balance,
		},
	}, nil
}

func (h *WalletHandler) GetTransactions(ctx context.Context, req *pb.GetTransactionsRequest) (*pb.GetTransactionsResponse, error) {
	userID := req.GetUserId()

	transactions, err := h.walletService.GetTransactions(ctx, userID)
	if err != nil {
		log.Printf("Failed to retrieve transactions for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to retrieve transactions: %v", err)
	}

	var protoTransactions []*pb.Transaction
	for _, tx := range transactions {
		protoTransactions = append(protoTransactions, &pb.Transaction{
			Id:     uint32(tx.TransactionID),
			UserId: uint32(tx.UserID),
			Type:   tx.TransactionType,
			Amount: float32(tx.Amount),
		})
	}

	return &pb.GetTransactionsResponse{
		Transactions: protoTransactions,
	}, nil
}
