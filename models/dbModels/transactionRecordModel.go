package dbModels

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionAction int

const (
	TransactionAction_NONE     TransactionAction = iota
	TransactionAction_Deposit                    // 入金
	TransactionAction_Withdraw                   // 出金
	TransactionAction_Bonus                      // 贈送
	TransactionAction_Interest                   // 利息
	TransactionAction_Open                       // 開倉
	TransactionAction_Close                      // 平倉
	TransactionAction_Manually                   // 人工更改
)

type TransactionStatus int

const (
	TransactionStatus_NONE     TransactionStatus = iota
	TransactionStatus_Pending                    // 待辦
	TransactionStatus_Success                    // 成功
	TransactionStatus_Failed                     // 失敗
	TransactionStatus_Rollback                   // 回滾
)

type TransactionRecordModel struct {
	ID                   uint64              `gorm:"column:id; primary_key"`
	MemberID             uint64              `gorm:"column:member_id"`
	WalletID             uint64              `gorm:"column:wallet_id"`
	Action               TransactionAction   `gorm:"column:action"`
	Amount               decimal.Decimal     `gorm:"column:amount"`
	BeforeAmount         decimal.NullDecimal `gorm:"column:before_amount"`
	AfterAmount          decimal.NullDecimal `gorm:"column:after_amount"`
	Currency             string              `gorm:"column:currency"`
	CommitterID          uint64              `gorm:"column:committer_id"`
	Status               TransactionStatus   `gorm:"column:status"`
	Remark               string              `gorm:"column:remark"`
	CreatedAt            time.Time           `gorm:"column:created_at"`
	UpdatedAt            time.Time           `gorm:"column:updated_at"`
	RollbackBeforeAmount decimal.NullDecimal `gorm:"column:rollback_before_amount"`
	RollbackAfterAmount  decimal.NullDecimal `gorm:"column:rollback_after_amount"`
	RollbackerID         sql.NullInt64       `gorm:"column:rollbacker_id"`
}
