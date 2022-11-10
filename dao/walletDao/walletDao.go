package walletDao

import (
	"errors"

	"github.com/paper-trade-chatbot/be-common/pagination"
	"github.com/paper-trade-chatbot/be-proto/general"
	"github.com/paper-trade-chatbot/be-wallet/models/dbModels"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

const table = "wallet"

// QueryModel set query condition, used by queryChain()
type QueryModel struct {
	ID       []uint64
	MemberID []uint64
	Currency *string
	Amount   *decimal.Decimal
}

type UpdateModel struct {
	Amount *decimal.Decimal
}

// New a row
func New(db *gorm.DB, model *dbModels.WalletModel) (int, error) {

	err := db.Table(table).
		Create(model).Error

	if err != nil {
		return 0, err
	}
	return 1, nil
}

// New rows
func News(db *gorm.DB, m []*dbModels.WalletModel) (int, error) {

	err := db.Transaction(func(tx *gorm.DB) error {

		err := tx.Table(table).
			CreateInBatches(m, 3000).Error

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return len(m), nil
}

// Get return a record as raw-data-form
func Get(tx *gorm.DB, query *QueryModel) (*dbModels.WalletModel, error) {

	result := &dbModels.WalletModel{}
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Scan(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Gets return records as raw-data-form
func Gets(tx *gorm.DB, query *QueryModel) ([]dbModels.WalletModel, error) {
	result := make([]dbModels.WalletModel, 0)
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Scan(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []dbModels.WalletModel{}, nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetsWithPagination(tx *gorm.DB, query *QueryModel, paginate *general.Pagination) ([]dbModels.WalletModel, *general.PaginationInfo, error) {

	var rows []dbModels.WalletModel
	var count int64 = 0
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Count(&count).
		Scopes(paginateChain(paginate)).
		Scan(&rows).Error

	offset, _ := pagination.GetOffsetAndLimit(paginate)
	paginationInfo := pagination.SetPaginationDto(paginate.Page, paginate.PageSize, int32(count), int32(offset))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []dbModels.WalletModel{}, paginationInfo, nil
	}

	if err != nil {
		return []dbModels.WalletModel{}, nil, err
	}

	return rows, paginationInfo, nil
}

// Modify a row
func Modify(tx *gorm.DB, model *dbModels.WalletModel, update *UpdateModel) error {
	attrs := map[string]interface{}{
		"amount": update.Amount,
	}

	err := tx.Table(table).
		Model(dbModels.WalletModel{}).
		Where(table+".id = ? AND "+table+".amount = ?", model.ID, model.Amount).
		Updates(attrs).Error

	return err
}

func Delete(db *gorm.DB, query *QueryModel) error {
	return db.Table(table).
		Scopes(queryChain(query)).
		Delete(&dbModels.WalletModel{}).Error
}

func queryChain(query *QueryModel) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Scopes(idInScope(query.ID)).
			Scopes(memberIDInScope(query.MemberID))

	}
}

func paginateChain(paginate *general.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset, limit := pagination.GetOffsetAndLimit(paginate)
		return db.
			Scopes(offsetScope(offset)).
			Scopes(limitScope(limit))

	}
}

func idInScope(id []uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(id) > 0 {
			return db.Where(table+".id IN ?", id)
		}
		return db
	}
}

func memberIDInScope(memberID []uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(memberID) > 0 {
			return db.Where(table+".member_id IN ?", memberID)
		}
		return db
	}
}

func currencyEqualScope(currency *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if currency != nil {
			return db.Where(table+".currency = ?", currency)
		}
		return db
	}
}

func limitScope(limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if limit > 0 {
			return db.Limit(limit)
		}
		return db
	}
}

func offsetScope(offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if offset > 0 {
			return db.Limit(offset)
		}
		return db
	}
}
