package coincache

import (
	"github.com/shopspring/decimal"
)

func float(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}
