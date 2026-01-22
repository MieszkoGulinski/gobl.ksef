package ksef_test

import (
	"testing"

	"github.com/invopop/gobl"
	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSign(t *testing.T) {
	createTestEnvelope := func() *gobl.Envelope {
		inv := &bill.Invoice{
			Currency:  currency.PLN,
			IssueDate: cal.MakeDate(2025, 12, 20),
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: l10n.PL.Tax(),
					Code:    "1234567890",
				},
			},
		}
		env, _ := gobl.Envelop(inv)
		return env
	}

	t.Run("returns error when envelope is nil", func(t *testing.T) {
		err := ksef.Sign(nil, "http://qr.url", "ksef-number", "hash")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "envelope is nil")
	})

	t.Run("adds KSeF number stamp", func(t *testing.T) {
		env := createTestEnvelope()

		err := ksef.Sign(env, "http://qr.url", "ksef-123456", "abc123hash")

		require.NoError(t, err)
		assert.True(t, hasStamp(env, favat.StampKSEFNumber, "ksef-123456"))
	})

	t.Run("adds hash stamp", func(t *testing.T) {
		env := createTestEnvelope()

		err := ksef.Sign(env, "http://qr.url", "ksef-123456", "abc123hash")

		require.NoError(t, err)
		assert.True(t, hasStamp(env, favat.StampHash, "abc123hash"))
	})

	t.Run("adds QR stamp", func(t *testing.T) {
		env := createTestEnvelope()

		err := ksef.Sign(env, "http://qr.example.com/verify", "ksef-123456", "abc123hash")

		require.NoError(t, err)
		assert.True(t, hasStamp(env, favat.StampQR, "http://qr.example.com/verify"))
	})

	t.Run("adds all three stamps", func(t *testing.T) {
		env := createTestEnvelope()

		err := ksef.Sign(env, "http://qr.url", "ksef-number", "hash-value")

		require.NoError(t, err)

		// Count stamps
		stampCount := 0
		for _, stamp := range env.Head.Stamps {
			if stamp.Provider == favat.StampKSEFNumber ||
				stamp.Provider == favat.StampHash ||
				stamp.Provider == favat.StampQR {
				stampCount++
			}
		}
		assert.Equal(t, 3, stampCount)
	})
}

func hasStamp(env *gobl.Envelope, provider cbc.Key, value string) bool {
	for _, stamp := range env.Head.Stamps {
		if stamp.Provider == provider && stamp.Value == value {
			return true
		}
	}
	return false
}
