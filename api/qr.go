package api

import (
	"encoding/base64"
	"fmt"
	"time"
)

const (
	EnvironmentProductionQrUrl = "qr.ksef.mf.gov.pl"
	EnvironmentDemoQrUrl       = "qr-demo.ksef.mf.gov.pl"
	EnvironmentTestQrUrl       = "qr-test.ksef.mf.gov.pl"
)

// GenerateQrCodeURL builds the verification URL for an invoice, both in online and offline mode.
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

	return fmt.Sprintf("https://%s/invoice/%s/%s/%s",
		baseUrl,
		nip,
		invoicingDate.Format("02-01-2006"),
		base64URLHash,
	), nil
}

// GenerateCertificateQrCodeURL builds the certificate verification URL for an offline invoice.
// The assembler url has format: base url / certificate / Nip / (contextNip) / (sellerNip) / (certificateId) / (invoiceHash)
// The url must not have a trailing slash or protocol.
func GenerateCertificateQrCodeURL(environment Environment, contextNip string, sellerNip string, certificateId string, invoiceHash []byte) (string, error) {
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

	if contextNip == "" {
		return "", fmt.Errorf("contextNip is empty")
	}
	if sellerNip == "" {
		return "", fmt.Errorf("sellerNip is empty")
	}
	if certificateId == "" {
		return "", fmt.Errorf("certificateId is empty")
	}

	base64URLHash := base64.RawURLEncoding.EncodeToString(invoiceHash)

	return fmt.Sprintf("%s/certificate/Nip/%s/%s/%s/%s",
		baseUrl,
		contextNip,
		sellerNip,
		certificateId,
		base64URLHash,
	), nil
}
