package main

import (
	"financial/entity"
	"financial/service"
	"fmt"
	"time"
)

func main() {
	// Initialize services
	userService := service.NewUserService()
	walletService := service.NewWalletService()
	recordService := service.NewRecordService()
	transactionCategoryService := service.NewTransactionCategoryService()
	reportService := service.NewReportService(walletService.GetWallets(), recordService.GetRecords(), transactionCategoryService.GetCategories())

	// Create a user
	user, _ := userService.CreateUser(&entity.User{Name: "John Doe", Email: "john@example.com", Password: "password123"})
	fmt.Println("Created User:", user)

	// Create wallets for the user
	wallet1, _ := walletService.CreateWallet(&entity.Wallet{UserID: user.ID, Name: "Personal Wallet", Balance: 1000.0})
	wallet2, _ := walletService.CreateWallet(&entity.Wallet{UserID: user.ID, Name: "Savings Wallet", Balance: 5000.0})
	fmt.Println("Created Wallets:", wallet1, wallet2)

	// Create transaction categories
	categoryFood, _ := transactionCategoryService.CreateTransactionCategory(&entity.TransactionCategory{Name: "Food"})
	categoryTransport, _ := transactionCategoryService.CreateTransactionCategory(&entity.TransactionCategory{Name: "Transport"})
	fmt.Println("Created Transaction Categories:", categoryFood, categoryTransport)

	// Create records
	record1, _ := recordService.CreateRecord(&entity.Record{
		WalletID:              wallet1.ID,
		TransactionCategoryID: categoryFood.ID,
		Amount:                50.0,
		Type:                  "expense",
		Timestamp:             time.Now().Format(time.RFC3339),
		Description:           "Lunch",
	})
	record2, _ := recordService.CreateRecord(&entity.Record{
		WalletID:              wallet2.ID,
		TransactionCategoryID: categoryTransport.ID,
		Amount:                20.0,
		Type:                  "expense",
		Timestamp:             time.Now().Format(time.RFC3339),
		Description:           "Taxi",
	})
	fmt.Println("Created Records:", record1, record2)

	// Transfer between wallets
	err := walletService.TransferBetweenWallets(wallet1.ID, wallet2.ID, 100.0)
	if err != nil {
		fmt.Println("Error transferring between wallets:", err)
	} else {
		fmt.Println("Transferred 100.0 from Personal Wallet to Savings Wallet")
	}

	// Get records by time range
	startTime := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	endTime := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
	records, _ := recordService.GetRecordsByTimeRange(wallet1.ID, startTime, endTime)
	fmt.Println("Records in Personal Wallet between time range:", records)

	// Get cash flow
	income, expense, _ := reportService.GetCashFlow(user.ID, startTime, endTime)
	fmt.Println("Cash Flow - Income:", income, "Expense:", expense)

	// Get expense recap by category
	expenseRecap, _ := reportService.GetExpenseRecapByCategory(user.ID, startTime, endTime)
	fmt.Println("Expense Recap by Category:", expenseRecap)

	// Get last 10 records
	last10Records, _ := recordService.GetLast10Records()
	fmt.Println("Last 10 Records:", last10Records)
}
