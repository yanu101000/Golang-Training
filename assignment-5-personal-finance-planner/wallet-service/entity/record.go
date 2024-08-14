package entity

type Record struct {
	ID                    int64   `json:"id"`
	WalletID              int64   `json:"wallet_id"`
	TransactionCategoryID int64   `json:"transaction_category_id"`
	Amount                float64 `json:"amount"`
	Type                  string  `json:"type"` // "income" or "expense"
	Timestamp             string  `json:"timestamp"`
	Description           string  `json:"description"`
}
