package ksef_test

import (
	"testing"

	"github.com/invopop/gobl.ksef/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildFAVAT(t *testing.T) {
	t.Run("should return a Document with KSeF data", func(t *testing.T) {
		doc, err := test.BuildFAVATFrom("invoice-standard.json")
		require.NoError(t, err)

		assert.Equal(t, "Faktura", doc.XMLName.Local)
		assert.Equal(t, "http://crd.gov.pl/wzor/2025/06/25/13775/", doc.XMLNamespace)
		assert.Equal(t, "http://www.w3.org/2001/XMLSchema", doc.XSDNamespace)
		assert.Equal(t, "http://www.w3.org/2001/XMLSchema-instance", doc.XSINamespace)
		assert.NotNil(t, doc.Header)
		assert.NotNil(t, doc.Buyer)
		assert.NotNil(t, doc.Seller)
		assert.NotNil(t, doc.Inv)
	})

	t.Run("should return bytes of the KSeF document", func(t *testing.T) {
		doc, err := test.BuildFAVATFrom("invoice-standard.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		output, err := test.LoadOutputFile("invoice-standard.xml")
		require.NoError(t, err)

		assert.Equal(t, string(output), string(data))
	})

	t.Run("should return bytes of the credit-note invoice", func(t *testing.T) {
		doc, err := test.BuildFAVATFrom("credit-note-standard.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		output, err := test.LoadOutputFile("credit-note-standard.xml")
		require.NoError(t, err)

		assert.Equal(t, string(output), string(data))
	})

	t.Run("should generate valid KSeF document", func(t *testing.T) {
		doc, err := test.BuildFAVATFrom("invoice-standard.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		test.ValidateAgainstFA3Schema(t, data)
	})

	t.Run("should generate valid credit-note", func(t *testing.T) {
		doc, err := test.BuildFAVATFrom("credit-note-standard.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		test.ValidateAgainstFA3Schema(t, data)
	})

	t.Run("should generate valid simplified invoice", func(t *testing.T) {
		doc, err := test.BuildFAVATFrom("invoice-simplified.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		test.ValidateAgainstFA3Schema(t, data)
	})

	t.Run("should generate valid self-billed invoice", func(t *testing.T) {
		doc, err := test.BuildFAVATFrom("invoice-self-billed.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		test.ValidateAgainstFA3Schema(t, data)

		output, err := test.LoadOutputFile("invoice-self-billed.xml")
		require.NoError(t, err)

		assert.Equal(t, string(output), string(data))
	})

	t.Run("should generate valid exempt invoice", func(t *testing.T) {
		doc, err := test.BuildFAVATFrom("invoice-exempt.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		test.ValidateAgainstFA3Schema(t, data)
	})

	t.Run("should generate valid reverse-charge invoice", func(t *testing.T) {
		doc, err := test.BuildFAVATFrom("invoice-reverse-charge.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		test.ValidateAgainstFA3Schema(t, data)
	})

	t.Run("should generate valid prepayment invoice", func(t *testing.T) {
		doc, err := test.BuildFAVATFrom("invoice-prepayment.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		test.ValidateAgainstFA3Schema(t, data)
	})

	t.Run("should generate valid settlement invoice", func(t *testing.T) {
		doc, err := test.BuildFAVATFrom("invoice-settlement.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		test.ValidateAgainstFA3Schema(t, data)
	})
}
