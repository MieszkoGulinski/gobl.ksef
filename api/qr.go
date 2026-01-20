package api

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// GenerateQrCodeURL builds the verification URL for an uploaded invoice.
func (c *Client) GenerateQrCodeURL(nip string, invoiceHash string, invoicingDate time.Time) (string, error) {
	if c == nil {
		return "", fmt.Errorf("client is nil")
	}
	if c.qrUrl == "" {
		return "", fmt.Errorf("qr url base is empty")
	}
	if nip == "" {
		return "", fmt.Errorf("nip is empty")
	}

	hashBytes, err := base64.StdEncoding.DecodeString(invoiceHash)
	if err != nil {
		return "", fmt.Errorf("invalid invoice hash: %w", err)
	}
	base64URLHash := base64.RawURLEncoding.EncodeToString(hashBytes)

	base := strings.TrimRight(c.qrUrl, "/")
	return fmt.Sprintf("%s/%s/%s/%s",
		base,
		nip,
		invoicingDate.Format("02-01-2006"),
		base64URLHash,
	), nil
}
