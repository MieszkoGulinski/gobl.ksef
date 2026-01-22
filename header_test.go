package ksef_test

import (
	"testing"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
)

func TestNewFavatHeader(t *testing.T) {
	t.Run("creates header with correct form code", func(t *testing.T) {
		inv := &bill.Invoice{
			IssueDate: cal.MakeDate(2025, 12, 20),
		}

		header := ksef.NewFavatHeader(inv)

		assert.NotNil(t, header.FormCode)
		assert.Equal(t, "FA (3)", header.FormCode.SystemCode)
		assert.Equal(t, "1-0E", header.FormCode.SchemaVersion)
		assert.Equal(t, "FA", header.FormCode.FormCode)
	})

	t.Run("sets form variant to 3", func(t *testing.T) {
		inv := &bill.Invoice{
			IssueDate: cal.MakeDate(2025, 12, 20),
		}

		header := ksef.NewFavatHeader(inv)

		assert.Equal(t, 3, header.FormVariant)
	})

	t.Run("formats creation date correctly", func(t *testing.T) {
		inv := &bill.Invoice{
			IssueDate: cal.MakeDate(2025, 6, 15),
		}

		header := ksef.NewFavatHeader(inv)

		assert.Equal(t, "2025-06-15T00:00:00Z", header.CreationDate)
	})

	t.Run("sets system info", func(t *testing.T) {
		inv := &bill.Invoice{
			IssueDate: cal.MakeDate(2025, 12, 20),
		}

		header := ksef.NewFavatHeader(inv)

		assert.Equal(t, "Invopop", header.SystemInfo)
	})
}
