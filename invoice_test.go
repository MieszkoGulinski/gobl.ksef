package ksef_test

import (
	"testing"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFavatInv(t *testing.T) {
	baseInvoice := func() *bill.Invoice {
		return &bill.Invoice{
			Currency: currency.PLN,
			Supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.PL.Tax(),
				},
			},
			Tax: &bill.Tax{
				Ext: tax.Extensions{
					favat.ExtKeyInvoiceType: "VAT",
				},
			},
			Totals: &bill.Totals{
				Taxes: &tax.Total{},
			},
		}
	}

	t.Run("sets preceding invoice", func(t *testing.T) {
		inv := baseInvoice()
		inv.Preceding = []*org.DocumentRef{
			{},
		}

		invoice := ksef.NewFavatInv(inv)

		assert.NotNil(t, invoice.CorrectedInv)
	})

	t.Run("sets correction reason", func(t *testing.T) {
		reason := "example reason"

		inv := baseInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				Reason: reason,
			},
		}

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, reason, invoice.CorrectionReason)
	})

	t.Run("sets correction type", func(t *testing.T) {
		inv := baseInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				Ext: tax.Extensions{
					favat.ExtKeyEffectiveDate: "1",
				},
			},
		}

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.CorrectionType)
	})

	t.Run("sets the self-billing annotation to false in non-self-billed invoices", func(t *testing.T) {
		inv := baseInvoice()

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "2", invoice.Annotations.SelfBilling)
	})

	t.Run("sets the self-billing annotation to true in self-billed invoices", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeySelfBilling] = "1"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.SelfBilling)
	})

	t.Run("sets reverse charge annotation", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyReverseCharge] = "1"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.ReverseCharge)
	})

	t.Run("sets cash accounting annotation", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyCashAccounting] = "1"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.CashAccounting)
	})

	t.Run("sets split payment annotation", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeySplitPayment] = "1"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.SplitPaymentMechanism)
	})

	t.Run("sets tax exemption annotation with marker", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyExemption] = "A"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.TaxExemption.Marker)
	})

	t.Run("sets margin scheme travel agency", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyMarginScheme] = "2"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.MarginScheme.Marker)
		assert.Equal(t, "1", invoice.Annotations.MarginScheme.TravelAgencyMargin)
	})

	t.Run("sets margin scheme used goods", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyMarginScheme] = "3.1"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.MarginScheme.Marker)
		assert.Equal(t, "1", invoice.Annotations.MarginScheme.UsedGoodsMargin)
	})

	t.Run("sets margin scheme art works", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyMarginScheme] = "3.2"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.MarginScheme.Marker)
		assert.Equal(t, "1", invoice.Annotations.MarginScheme.ArtWorksMargin)
	})

	t.Run("sets margin scheme collectibles and antiques", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyMarginScheme] = "3.3"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.MarginScheme.Marker)
		assert.Equal(t, "1", invoice.Annotations.MarginScheme.CollectiblesAndAntiquesMargin)
	})

	t.Run("sets additional description from notes", func(t *testing.T) {
		inv := baseInvoice()
		inv.Notes = []*org.Note{
			{
				Key:  "general",
				Text: "Test note text",
			},
		}

		invoice := ksef.NewFavatInv(inv)

		assert.Len(t, invoice.AdditionalDescription, 1)
		assert.Equal(t, "general", invoice.AdditionalDescription[0].Key)
		assert.Equal(t, "Test note text", invoice.AdditionalDescription[0].Value)
	})

	t.Run("sets invoice number with series", func(t *testing.T) {
		inv := baseInvoice()
		inv.Series = "INV"
		inv.Code = "001"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "INV-001", invoice.SequentialNumber)
	})

	t.Run("sets invoice number without series", func(t *testing.T) {
		inv := baseInvoice()
		inv.Code = "001"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "001", invoice.SequentialNumber)
	})

	t.Run("sets total amount due when due is specified", func(t *testing.T) {
		inv := baseInvoice()
		due, err := num.AmountFromString("10.00")
		require.NoError(t, err)
		inv.Totals.Due = &due

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "10.00", invoice.TotalAmountDue)
	})

	t.Run("sets total amount due from payable when due is nil", func(t *testing.T) {
		inv := baseInvoice()
		payable, err := num.AmountFromString("25.00")
		require.NoError(t, err)
		inv.Totals.Payable = payable

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "25.00", invoice.TotalAmountDue)
	})
}
