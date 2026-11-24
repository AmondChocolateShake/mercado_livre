package order

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AmondChocolateShake/mercado_livre/mercado-1999/internal/item"
)

var (
	ErrNotFound  = errors.New("order not found")
	ErrInvalidID = errors.New("invalid order id")

	errNotImplemented = errors.New("not implemented")
)

type CreateInput struct {
	ItemID   string
	BuyerID  int64
	Quantity int
}

type Service struct {
	db     *sql.DB
	orders *Repository
	items  *item.Repository
}

func NewService(db *sql.DB, orders *Repository, items *item.Repository) *Service {
	return &Service{db: db, orders: orders, items: items}
}

// TODO: 주문 생성 — 이 프로젝트의 핵심
// 규칙마다 위 var 블록에 ErrXxx 에러 변수를 정의해서 반환하세요
//
// 1. quantity > 0 검사
// 2. item.ParseID(in.ItemID)로 숫자 ID 얻기
// 3. database.WithTx(ctx, s.db, func(tx *sql.Tx) error { ... }) 안에서
//    a. s.items.WithTx(tx), s.orders.WithTx(tx)로 같은 트랜잭션을 쓰는 repository 얻기
//    b. FindByIDForUpdate로 상품 행 잠그고 읽기
//    c. 검사: 판매 중(active)인지 / 즉시구매(buy_it_now) 상품인지 / 판매자 본인 구매가 아닌지 / 재고가 충분한지
//    d. 재고 차감, 판매 수량 증가, 재고가 0이 되면 status를 closed로 → UpdateStock
//    e. Order 만들기 (unit_price, total_amount 계산, status=StatusConfirmed) → Insert
// 4. 트랜잭션 에러면 반환, 성공이면 생성된 주문 반환
//
// 생각해볼 것
// - FOR UPDATE를 빼면 재고 1개짜리 상품에 주문 2건이 동시에 들어올 때 어떻게 될까?
// - 재고 차감은 성공했는데 주문 Insert가 실패하면 재고는 어떻게 되어야 할까?
// - 트랜잭션 함수 안에서 만든 주문을 바깥으로 어떻게 꺼낼까?
func (s *Service) Create(ctx context.Context, in CreateInput) (*Order, error) {
	return nil, errNotImplemented
}

// TODO
// rawID를 숫자로 바꾸고 (실패하면 ErrInvalidID) s.orders.FindByID 호출
func (s *Service) Get(ctx context.Context, rawID string) (*Order, error) {
	return nil, errNotImplemented
}
