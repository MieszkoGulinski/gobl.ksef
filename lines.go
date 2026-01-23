package ksef

import (
	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Line defines the XML structure for KSeF item line (element type FaWiersz, for VAT and KOR type invoices)
type Line struct {
	LineNumber              int    `xml:"NrWierszaFa"`
	UniqueID                string `xml:"UU_ID,omitempty"`
	CompletionDate          string `xml:"P_6A,omitempty"`
	Name                    string `xml:"P_7,omitempty"`
	InternalCode            string `xml:"Indeks,omitempty"`
	GTIN                    string `xml:"GTIN,omitempty"`
	PKWiU                   string `xml:"PKWiU,omitempty"`
	CN                      string `xml:"CN,omitempty"`
	PKOB                    string `xml:"PKOB,omitempty"`
	Measure                 string `xml:"P_8A,omitempty"`
	Quantity                string `xml:"P_8B,omitempty"`
	NetUnitPrice            string `xml:"P_9A,omitempty"`
	GrossUnitPrice          string `xml:"P_9B,omitempty"`
	UnitDiscount            string `xml:"P_10,omitempty"`
	NetPriceTotal           string `xml:"P_11,omitempty"`
	GrossPriceTotal         string `xml:"P_11A,omitempty"`
	VATAmount               string `xml:"P_11Vat,omitempty"`
	VATRate                 string `xml:"P_12,omitempty"`
	OSSTaxRate              string `xml:"P_12_XII,omitempty"` // one stop shop
	Attachment15GoodsMarker int    `xml:"P_12_Zal_15,omitempty"`
	ExciseDuty              string `xml:"KwotaAkcyzy,omitempty"`
	SpecialGoodsCode        string `xml:"GTU,omitempty"` // values GTU_01 to GTU_13
	Procedure               string `xml:"Procedura,omitempty"`
	CurrencyRate            string `xml:"KursWaluty,omitempty"`
	BeforeCorrectionMarker  int    `xml:"StanPrzed,omitempty"`
}

// NewLines generates lines for the KSeF invoice
func NewLines(lines []*bill.Line) []*Line {
	var Lines []*Line

	for _, line := range lines {
		Lines = append(Lines, newLine(line))
	}

	return Lines
}

func newLine(line *bill.Line) *Line {
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
		if tc.Ext.Get(favat.ExtKeyTaxCategory) == "5" {
			if tc.Percent != nil {
				l.OSSTaxRate = tc.Percent.Amount().MinimalString()
			}
		} else {
			l.VATRate = vatRate(tc)
		}
	}

	return l
}

// vatRate returns the VAT rate string and OSS tax rate string for a tax combo
// based on the tax category extension
func vatRate(tc *tax.Combo) string {

	// For non-zero percentages, use the percentage value
	if tc.Percent != nil && !tc.Percent.IsZero() {
		return tc.Percent.Amount().MinimalString()
	}

	// For zero/nil percentage, determine from tax category extension
	switch tc.Ext.Get(favat.ExtKeyTaxCategory) {
	case "6.1": // zero-rated goods and services in the country
		return "0 KR"
	case "6.2": // intra-community supply
		return "0 WDT"
	case "6.3": // export supply
		return "0 EX"
	case "7": // tax exempt supply
		return "zw"
	case "8": // outside scope supply
		return "np I"
	case "9": // reverse charge supply
		return "np II"
	case "10": // domestic reverse charge supply
		return "oo"
	default:
		return ""
	}
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

// OrderLine defines the XML structure for KSeF item line (element type ZamowienieWiersz, for ZAL and KOR_ZAL type invoices)
type OrderLine struct {
	LineNumber              int    `xml:"NrWierszaZam"`
	UniqueID                string `xml:"UU_IDZ,omitempty"`
	Name                    string `xml:"P_7Z,omitempty"`
	InternalCode            string `xml:"IndeksZ,omitempty"`
	GTIN                    string `xml:"GTINZ,omitempty"`
	PKWiU                   string `xml:"PKWiUZ,omitempty"`
	CN                      string `xml:"CNZ,omitempty"`
	PKOB                    string `xml:"PKOBZ,omitempty"`
	Measure                 string `xml:"P_8AZ,omitempty"`
	Quantity                string `xml:"P_8BZ,omitempty"`
	NetUnitPrice            string `xml:"P_9AZ,omitempty"`
	NetPriceTotal           string `xml:"P_11NettoZ,omitempty"`
	TaxValue                string `xml:"P_11VatZ,omitempty"`
	VATRate                 string `xml:"P_12Z,omitempty"`
	OSSTaxRate              string `xml:"P_12Z_XII,omitempty"` // one stop shop
	Attachment15GoodsMarker int    `xml:"P_12Z_Zal_15,omitempty"`
	SpecialGoodsCode        string `xml:"GTUZ,omitempty"` // values GTU_01 to GTU_13
	Procedure               string `xml:"ProceduraZ,omitempty"`
	ExciseDuty              string `xml:"KwotaAkcyzyZ,omitempty"`
	BeforeCorrectionMarker  int    `xml:"StanPrzedZ,omitempty"`
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
