package item

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
)

const (
	maxTitleLength = 255
	maxSearchLimit = 50
)

var maxPrice = decimal.New(1, 13)

var (
	ErrNotFound          = errors.New("item not found")
	ErrInvalidID         = errors.New("invalid item id")
	ErrSellerNotFound    = errors.New("seller not found")
	ErrInvalidTitle      = errors.New("title must be 1 to 255 characters")
	ErrInvalidPrice      = errors.New("price must be greater than 0, below 10000000000000, with at most 2 decimal places")
	ErrInvalidQuantity   = errors.New("available_quantity must be greater than 0")
	ErrInvalidCondition  = errors.New("condition must be 'new' or 'used'")
	ErrInvalidBuyingMode = errors.New("buying_mode must be 'auction' or 'buy_it_now'")
	ErrAuctionQuantity   = errors.New("auction items must have available_quantity of 1")
	ErrInvalidStatus     = errors.New("status can only be changed to 'active' or 'paused'")
	ErrNoChanges         = errors.New("at least one field is required: title, price, available_quantity, status")
	ErrItemClosed        = errors.New("closed items cannot be modified")

	errNotImplemented = errors.New("not implemented")
)

type CreateInput struct {
	SellerID          int64
	Title             string
	Price             decimal.Decimal
	AvailableQuantity int
	Condition         string
	BuyingMode        string
}

type UpdateInput struct {
	Title             *string
	Price             *decimal.Decimal
	AvailableQuantity *int
	Status            *string
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*Item, error) {
	title, err := validateTitle(in.Title)
	if err != nil {
		return nil, err
	}
	if err := validatePrice(in.Price); err != nil {
		return nil, err
	}
	if in.AvailableQuantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if in.Condition != ConditionNew && in.Condition != ConditionUsed {
		return nil, ErrInvalidCondition
	}
	if in.BuyingMode != BuyingModeAuction && in.BuyingMode != BuyingModeBuyItNow {
		return nil, ErrInvalidBuyingMode
	}
	if in.BuyingMode == BuyingModeAuction && in.AvailableQuantity != 1 {
		return nil, ErrAuctionQuantity
	}

	it := &Item{
		SiteID:            DefaultSiteID,
		SellerID:          in.SellerID,
		Title:             title,
		Price:             in.Price,
		CurrencyID:        DefaultCurrencyID,
		AvailableQuantity: in.AvailableQuantity,
		BuyingMode:        in.BuyingMode,
		Condition:         in.Condition,
		Status:            StatusActive,
	}
	if err := s.repo.Insert(ctx, it); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, it.ID)
}

func (s *Service) Get(ctx context.Context, publicID string) (*Item, error) {
	id, err := ParseID(publicID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Search(ctx context.Context, keyword string, limit int) ([]Item, error) {
	if limit <= 0 || limit > maxSearchLimit {
		limit = maxSearchLimit
	}
	return s.repo.ListActive(ctx, strings.TrimSpace(keyword), limit)
}

func (s *Service) Update(ctx context.Context, publicID string, in UpdateInput) (*Item, error) {
	if in.Title == nil && in.Price == nil && in.AvailableQuantity == nil && in.Status == nil {
		return nil, ErrNoChanges
	}

	id, err := ParseID(publicID)
	if err != nil {
		return nil, err
	}

	it, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if it.Status == StatusClosed {
		return nil, ErrItemClosed
	}

	if in.Title != nil {
		title, err := validateTitle(*in.Title)
		if err != nil {
			return nil, err
		}
		it.Title = title
	}
	if in.Price != nil {
		if err := validatePrice(*in.Price); err != nil {
			return nil, err
		}
		it.Price = *in.Price
	}
	if in.AvailableQuantity != nil {
		if *in.AvailableQuantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		if it.BuyingMode == BuyingModeAuction && *in.AvailableQuantity != 1 {
			return nil, ErrAuctionQuantity
		}
		it.AvailableQuantity = *in.AvailableQuantity
	}
	if in.Status != nil {
		if *in.Status != StatusActive && *in.Status != StatusPaused {
			return nil, ErrInvalidStatus
		}
		it.Status = *in.Status
	}

	if err := s.repo.Update(ctx, it); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, publicID string) error {
	id, err := ParseID(publicID)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(ctx, id)
}

func validateTitle(raw string) (string, error) {
	title := strings.TrimSpace(raw)
	if title == "" || utf8.RuneCountInString(title) > maxTitleLength {
		return "", ErrInvalidTitle
	}
	return title, nil
}

func validatePrice(price decimal.Decimal) error {
	if !price.IsPositive() || price.GreaterThanOrEqual(maxPrice) || !price.Equal(price.Round(2)) {
		return ErrInvalidPrice
	}
	return nil
}
