package order

import (
	"context"
	"database/sql"

	"github.com/AmondChocolateShake/mercado_livre/mercado-1999/internal/database"
)

type Repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *Repository {
	return &Repository{db: db}
}

func (r *Repository) WithTx(tx *sql.Tx) *Repository {
	return &Repository{db: tx}
}

// TODO
// 1. INSERT INTO orders 실행 후 LastInsertId()로 o.ID 채우기
// 2. buyer_id가 users에 없으면 외래 키 에러 → 구매자 없음 에러로 바꿔 반환
func (r *Repository) Insert(ctx context.Context, o *Order) error {
	return errNotImplemented
}

// TODO
// 한 행 조회, 없으면 ErrNotFound
func (r *Repository) FindByID(ctx context.Context, id int64) (*Order, error) {
	return nil, errNotImplemented
}
