package item

import (
	"context"
	"errors"
)

var (
	ErrNotFound  = errors.New("item not found")
	ErrInvalidID = errors.New("invalid item id")

	errNotImplemented = errors.New("not implemented")
)

type CreateInput struct {
	SellerID          int64
	Title             string
	Price             int64
	AvailableQuantity int
	Condition         string
	BuyingMode        string
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// TODO: 상품 등록 규칙
// 규칙마다 위 var 블록에 ErrXxx 에러 변수를 정의해서 반환하세요
// 1. title 앞뒤 공백 제거 후 비어 있으면 에러
// 2. price > 0
// 3. available_quantity > 0
// 4. condition은 ConditionNew / ConditionUsed만
// 5. buying_mode는 BuyingModeAuction / BuyingModeBuyItNow만
// 6. 경매 상품은 수량 1개만 (왜 그래야 할까?)
// 7. Item 만들기: SiteID=DefaultSiteID, CurrencyID=DefaultCurrencyID, Status=StatusActive, DateCreated=현재 UTC 시각
// 8. s.repo.Insert 호출
func (s *Service) Create(ctx context.Context, in CreateInput) (*Item, error) {
	return nil, errNotImplemented
}

// TODO
// 1. ParseID로 "MLA100000000" → 100000000
// 2. s.repo.FindByID 호출
func (s *Service) Get(ctx context.Context, publicID string) (*Item, error) {
	return nil, errNotImplemented
}
