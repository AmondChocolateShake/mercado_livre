package order

import (
	"encoding/json"
	"time"

	"github.com/AmondChocolateShake/mercado_livre/mercado-1999/internal/item"
)

const StatusConfirmed = "confirmed"

type Order struct {
	ID          int64     `json:"id"`
	ItemID      int64     `json:"-"`
	BuyerID     int64     `json:"buyer_id"`
	SellerID    int64     `json:"seller_id"`
	Quantity    int       `json:"quantity"`
	UnitPrice   int64     `json:"unit_price"`
	TotalAmount int64     `json:"total_amount"`
	CurrencyID  string    `json:"currency_id"`
	Status      string    `json:"status"`
	DateCreated time.Time `json:"date_created"`
}

func (o Order) MarshalJSON() ([]byte, error) {
	type alias Order
	return json.Marshal(struct {
		alias
		ItemID string `json:"item_id"`
	}{
		alias:  alias(o),
		ItemID: item.FormatID(item.DefaultSiteID, o.ItemID),
	})
}
