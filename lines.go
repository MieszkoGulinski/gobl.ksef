package ksef

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Line defines the XML structure for KSeF item line (element type FaWiersz, for VAT and KOR type invoices)
type Line struct {
	LineNumber              int    `xml:"NrWierszaFa"`
	Name                    string `xml:"P_7,omitempty"`
	Measure                 string `xml:"P_8A,omitempty"`
	Quantity                string `xml:"P_8B,omitempty"`
	NetUnitPrice            string `xml:"P_9A,omitempty"`
	UnitDiscount            string `xml:"P_10,omitempty"`
	NetPriceTotal           string `xml:"P_11,omitempty"`
	VATRate                 string `xml:"P_12,omitempty"`
	ExciseDuty              string `xml:"KwotaAkcyzy,omitempty"`
	SpecialGoodsCode        string `xml:"GTU,omitempty"`      // values GTU_01 to GTU_13
	OSSTaxRate              string `xml:"P_12_XII,omitempty"` // one stop shop
	Attachment15GoodsMarker string `xml:"P_12_Zal_15,omitempty"`
	Procedure               string `xml:"Procedura,omitempty"`
	BeforeCorrectionMarker  string `xml:"StanPrzed,omitempty"`
}

func newLine(line *bill.Line, cu uint32) *Line {
	l := &Line{
		LineNumber:    line.Index,
		Name:          line.Item.Name,
		Measure:       string(line.Item.Unit.UNECE()),
		NetUnitPrice:  line.Item.Price.String(),
		Quantity:      line.Quantity.String(),
		UnitDiscount:  unitDiscount(line),
		NetPriceTotal: line.Total.String(),
	}
	if tc := line.Taxes.Get(tax.CategoryVAT); tc != nil {
		if tc.Percent != nil {
			l.VATRate = tc.Percent.Rescale(cu).StringWithoutSymbol()
		}
	}

	return l
}

func unitDiscount(line *bill.Line) string {
	if len(line.Discounts) == 0 {
		return ""
	}

	amount := num.MakeAmount(0, 2)

	for _, discount := range line.Discounts {
		amount = amount.Add(discount.Amount)
	}

	discount := amount.Divide(line.Quantity)

	return discount.String()
}

// NewLines generates lines for the KSeF invoice
func NewLines(lines []*bill.Line, cu uint32) []*Line {
	var Lines []*Line

	for _, line := range lines {
		Lines = append(Lines, newLine(line, cu))
	}

	return Lines
}

// OrderLine defines the XML structure for KSeF item line (element type ZamowienieWiersz, for ZAL and KOR_ZAL type invoices) - TODO use in the future
type OrderLine struct {
	LineNumber              int    `xml:"NrWierszaZam"`
	Name                    string `xml:"P_7Z,omitempty"`
	Measure                 string `xml:"P_8AZ,omitempty"`
	Quantity                string `xml:"P_8BZ,omitempty"`
	NetUnitPrice            string `xml:"P_9AZ,omitempty"`
	NetPriceTotal           string `xml:"P_11NettoZ,omitempty"`
	TaxValue                string `xml:"P_11VatZ,omitempty"`
	VATRate                 string `xml:"P_12Z,omitempty"`
	OSSTaxRate              string `xml:"P_12Z_XII,omitempty"` // one stop shop
	Attachment15GoodsMarker string `xml:"P_12Z_Zal_15,omitempty"`
	SpecialGoodsCode        string `xml:"GTUZ,omitempty"` // values GTU_1 to GTU_13
	Procedure               string `xml:"ProceduraZ,omitempty"`
	ExciseDuty              string `xml:"KwotaAkcyzyZ,omitempty"`
	BeforeCorrectionMarker  string `xml:"StanPrzedZ,omitempty"`
}

func newOrderLine(line *bill.Line, cu uint32) *OrderLine {
	l := &OrderLine{
		LineNumber:    line.Index,
		Name:          line.Item.Name,
		Measure:       string(line.Item.Unit.UNECE()),
		NetUnitPrice:  line.Item.Price.String(),
		Quantity:      line.Quantity.String(),
		NetPriceTotal: line.Total.String(),
	}
	if tc := line.Taxes.Get(tax.CategoryVAT); tc != nil {
		if tc.Percent != nil {
			l.VATRate = tc.Percent.Rescale(cu).StringWithoutSymbol()
			l.TaxValue = tc.Percent.Amount().Multiply(*line.Total).Rescale(cu).String() // TODO this is not correct
		}
	}

	return l
}

// NewOrderLines generates order lines for the KSeF invoice - TODO use in the future
func NewOrderLines(lines []*bill.Line, cu uint32) []*OrderLine {
	var orderLines []*OrderLine

	for _, line := range lines {
		orderLines = append(orderLines, newOrderLine(line, cu))
	}

	return orderLines
}
