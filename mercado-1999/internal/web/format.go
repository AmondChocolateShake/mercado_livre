package web

import (
	"html/template"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/AmondChocolateShake/mercado_livre/mercado-1999/internal/item"
)

var funcs = template.FuncMap{
	"price":      formatPrice,
	"date":       formatDate,
	"condition":  conditionLabel,
	"buyingMode": buyingModeLabel,
	"status":     statusLabel,
	"isAuction":  func(it item.Item) bool { return it.BuyingMode == item.BuyingModeAuction },
	"isActive":   func(it item.Item) bool { return it.Status == item.StatusActive },
}

func formatPrice(d decimal.Decimal) string {
	whole, cents, _ := strings.Cut(d.StringFixed(2), ".")

	var b strings.Builder
	for i, digit := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(digit)
	}

	if cents == "00" {
		return "$ " + b.String()
	}
	return "$ " + b.String() + "," + cents
}

func formatDate(t time.Time) string {
	return t.Format("02/01/2006")
}

func conditionLabel(condition string) string {
	switch condition {
	case item.ConditionNew:
		return "Nuevo"
	case item.ConditionUsed:
		return "Usado"
	}
	return condition
}

func buyingModeLabel(mode string) string {
	switch mode {
	case item.BuyingModeAuction:
		return "Subasta"
	case item.BuyingModeBuyItNow:
		return "Precio fijo"
	}
	return mode
}

func statusLabel(status string) string {
	switch status {
	case item.StatusActive:
		return "Activa"
	case item.StatusPaused:
		return "Pausada"
	case item.StatusClosed:
		return "Finalizada"
	}
	return status
}
