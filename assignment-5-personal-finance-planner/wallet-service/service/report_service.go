package service

import (
	"time"
	"wallet/entity"
)

type ReportService interface {
	GetCashFlow(userID int64, startTime, endTime string) (income float64, expense float64, err error)
	GetExpenseRecapByCategory(userID int64, startTime, endTime string) (map[string]float64, error)
}

type reportService struct {
	wallets    map[int64]*entity.Wallet
	records    map[int64]*entity.Record
	categories map[int64]*entity.TransactionCategory
}

func NewReportService(wallets map[int64]*entity.Wallet, records map[int64]*entity.Record, categories map[int64]*entity.TransactionCategory) ReportService {
	return &reportService{
		wallets:    wallets,
		records:    records,
		categories: categories,
	}
}

func (s *reportService) GetCashFlow(userID int64, startTime, endTime string) (income float64, expense float64, err error) {
	start, _ := time.Parse(time.RFC3339, startTime)
	end, _ := time.Parse(time.RFC3339, endTime)

	for _, wallet := range s.wallets {
		if wallet.UserID == userID {
			for _, record := range s.records {
				if record.WalletID == wallet.ID {
					recordTime, _ := time.Parse(time.RFC3339, record.Timestamp)
					if recordTime.After(start) && recordTime.Before(end) {
						if record.Type == "income" {
							income += record.Amount
						} else if record.Type == "expense" {
							expense += record.Amount
						}
					}
				}
			}
		}
	}

	return income, expense, nil
}

func (s *reportService) GetExpenseRecapByCategory(userID int64, startTime, endTime string) (map[string]float64, error) {
	start, _ := time.Parse(time.RFC3339, startTime)
	end, _ := time.Parse(time.RFC3339, endTime)
	recap := make(map[string]float64)

	for _, wallet := range s.wallets {
		if wallet.UserID == userID {
			for _, record := range s.records {
				if record.WalletID == wallet.ID {
					recordTime, _ := time.Parse(time.RFC3339, record.Timestamp)
					if recordTime.After(start) && recordTime.Before(end) && record.Type == "expense" {
						category := s.categories[record.TransactionCategoryID].Name
						recap[category] += record.Amount
					}
				}
			}
		}
	}

	return recap, nil
}
