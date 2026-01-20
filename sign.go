package ksef

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/regimes/pl"
)

// Sign attaches KSeF identification stamps (KSeF number, hash and QR URL) to the envelope.
func Sign(env *gobl.Envelope, qrURL, ksefNumber, invoiceHash string) error {
	if env == nil {
		return fmt.Errorf("envelope is nil")
	}

	env.Head.AddStamp(
		&head.Stamp{
			Provider: pl.StampProviderKSeFID,
			Value:    ksefNumber,
		},
	)
	env.Head.AddStamp(
		&head.Stamp{
			Provider: pl.StampProviderKSeFHash,
			Value:    invoiceHash,
		},
	)
	env.Head.AddStamp(
		&head.Stamp{
			Provider: pl.StampProviderKSeFQR,
			Value:    qrURL,
		},
	)

	return nil
}

// GenerateQrCodeURL builds the URL used for verifying an invoice in KSeF.
func GenerateQrCodeURL(baseURL, nip, invoiceHash string, invoicingDate time.Time) (string, error) {
	if baseURL == "" {
		return "", fmt.Errorf("base URL is empty")
	}
	if nip == "" {
		return "", fmt.Errorf("nip is empty")
	}

	hashBytes, err := base64.StdEncoding.DecodeString(invoiceHash)
	if err != nil {
		return "", fmt.Errorf("invalid invoice hash: %w", err)
	}
	base64URLHash := base64.RawURLEncoding.EncodeToString(hashBytes)

	base := strings.TrimRight(baseURL, "/") // remove trailing slash if present
	return fmt.Sprintf("%s/%s/%s/%s",
		base,
		nip,
		invoicingDate.Format("02-01-2006"),
		base64URLHash,
	), nil
}
