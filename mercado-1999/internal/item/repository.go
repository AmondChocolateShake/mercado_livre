package item

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
// 1. r.db.ExecContext로 INSERT INTO items 실행 (값은 ? 플레이스홀더로)
// 2. res.LastInsertId()로 받은 id를 it.ID에 채우기
// 3. seller_id가 users에 없으면 외래 키 에러가 난다
//    database.IsForeignKeyViolation(err)로 판별해서 판매자 없음 에러로 바꿔 반환
func (r *Repository) Insert(ctx context.Context, it *Item) error {
	return errNotImplemented
}

// TODO
// 1. r.db.QueryRowContext(...).Scan(...)으로 한 행 읽기
// 2. 행이 없으면 sql.ErrNoRows가 나온다 → ErrNotFound로 바꿔 반환
func (r *Repository) FindByID(ctx context.Context, id int64) (*Item, error) {
	return nil, errNotImplemented
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
