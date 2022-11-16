package transactionRecordDao

import (
	"database/sql"
	"errors"
	"time"

	"github.com/paper-trade-chatbot/be-common/pagination"
	"github.com/paper-trade-chatbot/be-proto/general"
	"github.com/paper-trade-chatbot/be-wallet/models/dbModels"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

const table = "transaction_record"

type OrderColumn int

const (
	OrderColumn_None OrderColumn = iota
	OrderColumn_MemberID
	OrderColumn_CommitterID
	OrderColumn_Currency
	OrderColumn_CreatedAt
)

type OrderDirection int

const (
	OrderDirection_None = 0
	OrderDirection_ASC  = 1
	OrderDirection_DESC = -1
)

type Order struct {
	Column    OrderColumn
	Direction OrderDirection
}

// QueryModel set query condition, used by queryChain()
type QueryModel struct {
	ID           *uint64
	MemberID     *uint64
	CommitterID  *uint64
	RollbackerID *uint64
	Currency     []string
	Action       []dbModels.TransactionAction
	Status       []dbModels.TransactionStatus
	CreatedFrom  *time.Time
	CreatedTo    *time.Time
	OrderBy      []*Order
}

type UpdateModel struct {
	BeforeAmount         *decimal.NullDecimal
	AfterAmount          *decimal.NullDecimal
	Status               *dbModels.TransactionStatus
	Remark               *sql.NullString
	RollbackBeforeAmount *decimal.NullDecimal
	RollbackAfterAmount  *decimal.NullDecimal
	RollbackerID         *sql.NullInt64
}

// New a row
func New(db *gorm.DB, model *dbModels.TransactionRecordModel) (int, error) {

	err := db.Table(table).
		Create(model).Error

	if err != nil {
		return 0, err
	}
	return 1, nil
}

// New rows
func News(db *gorm.DB, m []*dbModels.TransactionRecordModel) (int, error) {

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
func Get(tx *gorm.DB, query *QueryModel) (*dbModels.TransactionRecordModel, error) {

	result := &dbModels.TransactionRecordModel{}
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
func Gets(tx *gorm.DB, query *QueryModel) ([]dbModels.TransactionRecordModel, error) {
	result := make([]dbModels.TransactionRecordModel, 0)
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Scan(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []dbModels.TransactionRecordModel{}, nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetsWithPagination(tx *gorm.DB, query *QueryModel, paginate *general.Pagination) ([]dbModels.TransactionRecordModel, *general.PaginationInfo, error) {

	var rows []dbModels.TransactionRecordModel
	var count int64 = 0
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Count(&count).
		Scopes(paginateChain(paginate)).
		Scan(&rows).Error

	offset, _ := pagination.GetOffsetAndLimit(paginate)
	paginationInfo := pagination.SetPaginationDto(paginate.Page, paginate.PageSize, int32(count), int32(offset))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []dbModels.TransactionRecordModel{}, paginationInfo, nil
	}

	if err != nil {
		return []dbModels.TransactionRecordModel{}, nil, err
	}

	return rows, paginationInfo, nil
}

// Gets return records as raw-data-form
func Modify(tx *gorm.DB, model *dbModels.TransactionRecordModel, update *UpdateModel) error {
	attrs := map[string]interface{}{}
	if update.BeforeAmount != nil {
		attrs["before_amount"] = *update.BeforeAmount
	}
	if update.AfterAmount != nil {
		attrs["after_amount"] = *update.AfterAmount
	}
	if update.Status != nil {
		attrs["status"] = *update.Status
	}
	if update.Remark != nil {
		attrs["remark"] = *update.Remark
	}
	if update.RollbackBeforeAmount != nil {
		attrs["rollback_before_amount"] = *update.RollbackBeforeAmount
	}
	if update.RollbackAfterAmount != nil {
		attrs["rollback_after_amount"] = *update.RollbackAfterAmount
	}
	if update.RollbackerID != nil {
		attrs["rollbacker_id"] = *update.RollbackerID
	}

	err := tx.Table(table).
		Model(dbModels.TransactionRecordModel{}).
		Where(table+".id = ? AND "+table+".status = ?", model.ID, model.Status).
		Updates(attrs).Error

	return err
}

func queryChain(query *QueryModel) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Scopes(idEqualScope(query.ID)).
			Scopes(memberIDEqualScope(query.MemberID)).
			Scopes(committerIDEqualScope(query.CommitterID)).
			Scopes(rollbackerIDEqualScope(query.RollbackerID)).
			Scopes(currencyInScope(query.Currency)).
			Scopes(statusInScope(query.Status)).
			Scopes(actionInScope(query.Action)).
			Scopes(createdBetweenScope(query.CreatedFrom, query.CreatedTo)).
			Scopes(orderByScope(query.OrderBy))
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

func idEqualScope(id *uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if id != nil {
			return db.Where(table+".id = ?", *id)
		}
		return db
	}
}

func memberIDEqualScope(memberID *uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if memberID != nil {
			return db.Where(table+".member_id = ?", *memberID)
		}
		return db
	}
}

func committerIDEqualScope(committerID *uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if committerID != nil {
			return db.Where(table+".committer_id = ?", *committerID)
		}
		return db
	}
}

func rollbackerIDEqualScope(committerID *uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if committerID != nil {
			return db.Where(table+".committer_id = ?", *committerID)
		}
		return db
	}
}

func currencyInScope(currency []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(currency) > 0 {
			return db.Where(table+".currency IN ?", currency)
		}
		return db
	}
}

func statusInScope(status []dbModels.TransactionStatus) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(status) > 0 {
			return db.Where(table+".status IN ?", status)
		}
		return db
	}
}

func actionInScope(action []dbModels.TransactionAction) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(action) > 0 {
			return db.Where(table+".action IN ?", action)
		}
		return db
	}
}

func createdBetweenScope(createdFrom, createdTo *time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if createdFrom != nil && createdTo != nil {
			return db.Where(table+".created_at BETWEEN ? AND ?", createdFrom, createdTo)
		}
		return db
	}
}

func orderByScope(order []*Order) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(order) > 0 {
			for _, o := range order {
				orderClause := ""
				switch o.Column {
				case OrderColumn_MemberID:
					orderClause += "member_id"
				case OrderColumn_CommitterID:
					orderClause += "committer_id"
				case OrderColumn_Currency:
					orderClause += "currency"
				case OrderColumn_CreatedAt:
					orderClause += "created_at"
				default:
					continue
				}

				switch o.Direction {
				case OrderDirection_ASC:
					orderClause += " ASC"
				case OrderDirection_DESC:
					orderClause += " DESC"
				}

				db = db.Order(orderClause)
			}
			return db
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
