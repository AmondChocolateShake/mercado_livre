package item

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultSiteID     = "MLA"
	DefaultCurrencyID = "ARS"

	ConditionNew  = "new"
	ConditionUsed = "used"

	BuyingModeAuction  = "auction"
	BuyingModeBuyItNow = "buy_it_now"

	StatusActive = "active"
	StatusPaused = "paused"
	StatusClosed = "closed"
)

type Item struct {
	ID                int64     `json:"-"`
	SiteID            string    `json:"site_id"`
	SellerID          int64     `json:"seller_id"`
	Title             string    `json:"title"`
	Price             int64     `json:"price"`
	CurrencyID        string    `json:"currency_id"`
	AvailableQuantity int       `json:"available_quantity"`
	SoldQuantity      int       `json:"sold_quantity"`
	BuyingMode        string    `json:"buying_mode"`
	Condition         string    `json:"condition"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

func (i Item) PublicID() string {
	return FormatID(i.SiteID, i.ID)
}

func (i Item) MarshalJSON() ([]byte, error) {
	type alias Item
	return json.Marshal(struct {
		ID string `json:"id"`
		alias
	}{
		ID:    i.PublicID(),
		alias: alias(i),
	})
}

func FormatID(siteID string, id int64) string {
	return siteID + strconv.FormatInt(id, 10)
}

func ParseID(publicID string) (int64, error) {
	numeric, ok := strings.CutPrefix(publicID, DefaultSiteID)
	if !ok {
		return 0, ErrInvalidID
	}

	id, err := strconv.ParseInt(numeric, 10, 64)
	if err != nil || id <= 0 {
		return 0, ErrInvalidID
	}
	return id, nil
}
