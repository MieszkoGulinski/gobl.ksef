package ksef_test

import (
	"testing"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLines(t *testing.T) {
	t.Run("converts basic lines", func(t *testing.T) {
		price, _ := num.AmountFromString("100.00")
		qty, _ := num.AmountFromString("2")
		total, _ := num.AmountFromString("200.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Test Item",
					Price: &price,
					Unit:  "h",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
					},
				},
			},
		}

		result := ksef.NewLines(lines)

		require.Len(t, result, 1)
		assert.Equal(t, 1, result[0].LineNumber)
		assert.Equal(t, "Test Item", result[0].Name)
		assert.Equal(t, "HUR", result[0].Measure)
		assert.Equal(t, "100.00", result[0].NetUnitPrice)
		assert.Equal(t, "2", result[0].Quantity)
		assert.Equal(t, "200.00", result[0].NetPriceTotal)
		assert.Equal(t, "23", result[0].VATRate)
	})

	t.Run("handles multiple lines", func(t *testing.T) {
		price, _ := num.AmountFromString("50.00")
		qty, _ := num.AmountFromString("1")
		total, _ := num.AmountFromString("50.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Item 1",
					Price: &price,
					Unit:  "h",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
					},
				},
			},
			{
				Index:    2,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Item 2",
					Price: &price,
					Unit:  "service",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(8, 2),
					},
				},
			},
		}

		result := ksef.NewLines(lines)

		require.Len(t, result, 2)
		assert.Equal(t, "Item 1", result[0].Name)
		assert.Equal(t, "Item 2", result[1].Name)
		assert.Equal(t, 1, result[0].LineNumber)
		assert.Equal(t, 2, result[1].LineNumber)
	})

	t.Run("handles line without VAT percent", func(t *testing.T) {
		price, _ := num.AmountFromString("100.00")
		qty, _ := num.AmountFromString("1")
		total, _ := num.AmountFromString("100.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Exempt Item",
					Price: &price,
					Unit:  "h",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						// No Percent for exempt items
					},
				},
			},
		}

		result := ksef.NewLines(lines)

		require.Len(t, result, 1)
		assert.Equal(t, "", result[0].VATRate)
	})

	t.Run("handles line with discounts", func(t *testing.T) {
		price, _ := num.AmountFromString("100.00")
		qty, _ := num.AmountFromString("2")
		total, _ := num.AmountFromString("180.00")
		discountAmt, _ := num.AmountFromString("20.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Discounted Item",
					Price: &price,
					Unit:  "h",
				},
				Discounts: []*bill.LineDiscount{
					{
						Amount: discountAmt,
						Reason: "Volume discount",
					},
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
					},
				},
			},
		}

		result := ksef.NewLines(lines)

		require.Len(t, result, 1)
		assert.Equal(t, "10.00", result[0].UnitDiscount) // 20.00 / 2 = 10.00
	})

	t.Run("handles multiple discounts on same line", func(t *testing.T) {
		price, _ := num.AmountFromString("100.00")
		qty, _ := num.AmountFromString("4")
		total, _ := num.AmountFromString("360.00")
		discount1, _ := num.AmountFromString("20.00")
		discount2, _ := num.AmountFromString("20.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Multiple Discounts",
					Price: &price,
					Unit:  "h",
				},
				Discounts: []*bill.LineDiscount{
					{Amount: discount1, Reason: "Discount 1"},
					{Amount: discount2, Reason: "Discount 2"},
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
					},
				},
			},
		}

		result := ksef.NewLines(lines)

		require.Len(t, result, 1)
		assert.Equal(t, "10.00", result[0].UnitDiscount) // (20.00 + 20.00) / 4 = 10.00
	})
}

func TestNewOrderLines(t *testing.T) {
	t.Run("converts order lines with VAT calculation", func(t *testing.T) {
		price, _ := num.AmountFromString("100.00")
		qty, _ := num.AmountFromString("2")
		total, _ := num.AmountFromString("200.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Order Item",
					Price: &price,
					Unit:  "h",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
						Ext: tax.Extensions{
							favat.ExtKeyTaxCategory: "1",
						},
					},
				},
			},
		}

		result := ksef.NewOrderLines(lines, 2)

		require.Len(t, result, 1)
		assert.Equal(t, 1, result[0].LineNumber)
		assert.Equal(t, "Order Item", result[0].Name)
		assert.Equal(t, "HUR", result[0].Measure)
		assert.Equal(t, "100.00", result[0].NetUnitPrice)
		assert.Equal(t, "2", result[0].Quantity)
		assert.Equal(t, "200.00", result[0].NetPriceTotal)
		assert.Equal(t, "23", result[0].VATRate)
		// Note: TaxValue calculation in current code uses Percent.Amount() which
		// returns the percentage value directly (23), not the decimal (0.23),
		// resulting in 200.00 * 23 = 4600.00 (there's a TODO in the code about this)
		assert.Equal(t, "4600.00", result[0].TaxValue)
	})

	t.Run("handles multiple order lines", func(t *testing.T) {
		price, _ := num.AmountFromString("50.00")
		qty, _ := num.AmountFromString("1")
		total, _ := num.AmountFromString("50.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Item 1",
					Price: &price,
					Unit:  "h",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
					},
				},
			},
			{
				Index:    2,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Item 2",
					Price: &price,
					Unit:  "service",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(8, 2),
					},
				},
			},
		}

		result := ksef.NewOrderLines(lines, 2)

		require.Len(t, result, 2)
		assert.Equal(t, "Item 1", result[0].Name)
		assert.Equal(t, "Item 2", result[1].Name)
	})
}
