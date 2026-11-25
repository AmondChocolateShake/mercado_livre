package item

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AmondChocolateShake/mercado_livre/mercado-1999/internal/database"
)

const selectColumns = `id, site_id, seller_id, title, price, currency_id,
	available_quantity, sold_quantity, buying_mode, item_condition, status,
	created_at, updated_at, deleted_at`

type Repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *Repository {
	return &Repository{db: db}
}

func (r *Repository) WithTx(tx *sql.Tx) *Repository {
	return &Repository{db: tx}
}

func (r *Repository) Insert(ctx context.Context, it *Item) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO items (site_id, seller_id, title, price, currency_id,
			available_quantity, sold_quantity, buying_mode, item_condition, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		it.SiteID, it.SellerID, it.Title, it.Price, it.CurrencyID,
		it.AvailableQuantity, it.SoldQuantity, it.BuyingMode, it.Condition, it.Status,
	)
	if err != nil {
		if database.IsForeignKeyViolation(err) {
			return ErrSellerNotFound
		}
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	it.ID = id
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*Item, error) {
	var it Item
	err := r.db.QueryRowContext(ctx,
		`SELECT `+selectColumns+` FROM items WHERE id = ? AND deleted_at IS NULL`, id,
	).Scan(
		&it.ID, &it.SiteID, &it.SellerID, &it.Title, &it.Price, &it.CurrencyID,
		&it.AvailableQuantity, &it.SoldQuantity, &it.BuyingMode, &it.Condition, &it.Status,
		&it.CreatedAt, &it.UpdatedAt, &it.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &it, nil
}

func (r *Repository) Update(ctx context.Context, it *Item) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE items
		SET title = ?, price = ?, available_quantity = ?, status = ?
		WHERE id = ? AND deleted_at IS NULL AND status <> ?`,
		it.Title, it.Price, it.AvailableQuantity, it.Status,
		it.ID, StatusClosed,
	)
	if err != nil {
		return err
	}
	return expectOneRow(res, ErrItemClosed)
}

func (r *Repository) SoftDelete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE items
		SET deleted_at = UTC_TIMESTAMP(), status = ?
		WHERE id = ? AND deleted_at IS NULL`,
		StatusClosed, id,
	)
	if err != nil {
		return err
	}
	return expectOneRow(res, ErrNotFound)
}

// TODO
// FindByID와 같은 쿼리 끝에 FOR UPDATE를 붙인다
// 생각해볼 것: 트랜잭션 밖에서 FOR UPDATE를 쓰면 어떻게 될까?
func (r *Repository) FindByIDForUpdate(ctx context.Context, id int64) (*Item, error) {
	return nil, errNotImplemented
}

// TODO
// available_quantity, sold_quantity, status를 UPDATE
func (r *Repository) UpdateStock(ctx context.Context, id int64, available, sold int, status string) error {
	return errNotImplemented
}

func expectOneRow(res sql.Result, errIfNone error) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errIfNone
	}
	return nil
}
