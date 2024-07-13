package entity

type Wallet struct {
	WalletID uint    `gorm:"primaryKey"`
	UserID   uint    `gorm:"not null"`
	Balance  float32 `gorm:"type:numeric(10,2);default:0"`
}
