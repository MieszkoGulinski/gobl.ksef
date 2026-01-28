package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateQrCodeURL(t *testing.T) {
	invoiceHash := []byte("hash1234567890") // Example hash

	t.Run("production", func(t *testing.T) {
		date, _ := time.Parse("2006-01-02", "2023-10-25")
		url, err := GenerateQrCodeURL(EnvironmentProduction, "1234567890", date, invoiceHash)
		assert.NoError(t, err)
		assert.Equal(t, "https://qr.ksef.mf.gov.pl/invoice/1234567890/25-10-2023/aGFzaDEyMzQ1Njc4OTA", url)
	})

	t.Run("test", func(t *testing.T) {
		date, _ := time.Parse("2006-01-02", "2023-10-25")
		url, err := GenerateQrCodeURL(EnvironmentTest, "1234567890", date, invoiceHash)
		assert.NoError(t, err)
		assert.Equal(t, "https://qr-test.ksef.mf.gov.pl/invoice/1234567890/25-10-2023/aGFzaDEyMzQ1Njc4OTA", url)
	})
}

func TestGenerateCertificateQrCodeURL(t *testing.T) {
	invoiceHash := []byte{0x00, 0x01, 0x02, 0x03}
	// Base64 RawURL of 00010203 -> AAECAw

	t.Run("test environment", func(t *testing.T) {
		url, err := GenerateCertificateQrCodeURL(
			EnvironmentTest,
			"1111111111", // contextNip
			"2222222222", // sellerNip
			"CERT123",    // certificateId
			invoiceHash,
		)
		assert.NoError(t, err)
		assert.Equal(t, "qr-test.ksef.mf.gov.pl/certificate/Nip/1111111111/2222222222/CERT123/AAECAw", url)
	})

	t.Run("production environment", func(t *testing.T) {
		url, err := GenerateCertificateQrCodeURL(
			EnvironmentProduction,
			"1111111111",
			"2222222222",
			"CERT123",
			invoiceHash,
		)
		assert.NoError(t, err)
		assert.Equal(t, "qr.ksef.mf.gov.pl/certificate/Nip/1111111111/2222222222/CERT123/AAECAw", url)
	})

	t.Run("empty args", func(t *testing.T) {
		_, err := GenerateCertificateQrCodeURL(EnvironmentTest, "", "222", "C", invoiceHash)
		assert.Error(t, err)
		_, err = GenerateCertificateQrCodeURL(EnvironmentTest, "111", "", "C", invoiceHash)
		assert.Error(t, err)
		_, err = GenerateCertificateQrCodeURL(EnvironmentTest, "111", "222", "", invoiceHash)
		assert.Error(t, err)
	})
}
