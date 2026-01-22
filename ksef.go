// Package ksef implements conversion between GOBL documents and KSeF formats,
// including the Polish FA_VAT XML invoice document.
package ksef

import (
	"encoding/xml"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
)

// Constants for KSeF XML
const (
	XSINamespace    = "http://www.w3.org/2001/XMLSchema-instance"
	XSDNamespace    = "http://www.w3.org/2001/XMLSchema"
	XMLNamespace    = "http://crd.gov.pl/wzor/2025/06/25/13775/"
	RootElementName = "Faktura"
)

// Invoice is a pseudo-model for containing the XML document being created
type Invoice struct {
	XMLName      xml.Name
	XSINamespace string        `xml:"xmlns:xsi,attr"`
	XSDNamespace string        `xml:"xmlns:xsd,attr"`
	XMLNamespace string        `xml:"xmlns,attr"`
	Header       *Header       `xml:"Naglowek"`
	Seller       *Seller       `xml:"Podmiot1"`
	Buyer        *Buyer        `xml:"Podmiot2"`
	ThirdParties []*ThirdParty `xml:"Podmiot3,omitempty"` // third party (up to 100)
	Inv          *Inv          `xml:"Fa"`
}

// BuildFavat converts a GOBL envelope into a KSeF FA_VAT invoice document.
func BuildFavat(env *gobl.Envelope) (*Invoice, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}

	if !favat.V3.In(inv.GetAddons()...) {
		return nil, fmt.Errorf("invoice does not have the FA_VAT v3 addon")
	}

	if inv.Type == bill.InvoiceTypeCreditNote {
		// In KSEF credit notes become corrective invoices,
		// which require negative totals.
		if err := inv.Invert(); err != nil {
			return nil, err
		}
	}

	invoice := &Invoice{
		XMLName:      xml.Name{Local: RootElementName},
		XSINamespace: XSINamespace,
		XSDNamespace: XSDNamespace,
		XMLNamespace: XMLNamespace,

		Header:       NewFavatHeader(inv),
		Seller:       NewFavatSeller(inv.Supplier),
		Buyer:        NewFavatBuyer(inv.Customer),
		ThirdParties: NewThirdParties(inv),
		Inv:          NewFavatInv(inv),
	}

	return invoice, nil
}

// Bytes returns the XML representation of the document in bytes
func (d *Invoice) Bytes() ([]byte, error) {
	data, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), data...), nil
}
