package dbModels

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type WalletModel struct {
	ID        uint64          `gorm:"column:id; primary_key"`
	MemberID  uint64          `gorm:"column:member_id"`
	Amount    decimal.Decimal `gorm:"column:amount"`
	Currency  string          `gorm:"column:currency"`
	CreatedAt *time.Time      `gorm:"column:created_at"`
	UpdatedAt *time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at"`
}
