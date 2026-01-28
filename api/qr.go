package api

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

const (
	EnvironmentProductionQrUrl = "https://qr.ksef.mf.gov.pl/invoice"
	EnvironmentDemoQrUrl       = "https://qr-demo.ksef.mf.gov.pl/invoice"
	EnvironmentTestQrUrl       = "https://qr-test.ksef.mf.gov.pl/invoice"
)

// GenerateQrCodeURL builds the verification URL for an uploaded invoice.
func GenerateQrCodeURL(environment Environment, nip string, invoicingDate time.Time, invoiceHash []byte) (string, error) {
	var baseUrl string
	switch environment {
	case EnvironmentProduction:
		baseUrl = EnvironmentProductionQrUrl
	case EnvironmentDemo:
		baseUrl = EnvironmentDemoQrUrl
	case EnvironmentTest:
		baseUrl = EnvironmentTestQrUrl
	default:
		return "", fmt.Errorf("invalid environment: %s", environment)
	}

	if nip == "" {
		return "", fmt.Errorf("nip is empty")
	}

	base64URLHash := base64.RawURLEncoding.EncodeToString(invoiceHash)

	base := strings.TrimRight(baseUrl, "/")
	return fmt.Sprintf("%s/%s/%s/%s",
		base,
		nip,
		invoicingDate.Format("02-01-2006"),
		base64URLHash,
	), nil
}
