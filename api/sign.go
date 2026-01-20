package api

import (
	"encoding/base64"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/regimes/pl"
)

// Sign attached QR code and other identification values to the envelope
func (c *Client) Sign(env *gobl.Envelope, nip string, uploadedInvoice *UploadedInvoice) error {
	env.Head.AddStamp(
		&head.Stamp{
			Provider: pl.StampProviderKSeFID,
			Value:    uploadedInvoice.KsefNumber,
		},
	)
	env.Head.AddStamp(
		&head.Stamp{
			Provider: pl.StampProviderKSeFHash,
			Value:    uploadedInvoice.InvoiceHash,
		},
	)
	// URL contains invoicing date in DD-MM-YYYY format
	// Hash must be in Base64URL, not Base64
	hashBytes, err := base64.StdEncoding.DecodeString(uploadedInvoice.InvoiceHash)
	if err != nil {
		return err
	}
	base64UrlHash := base64.RawURLEncoding.EncodeToString(hashBytes)

	env.Head.AddStamp(
		&head.Stamp{
			Provider: pl.StampProviderKSeFQR,
			Value:    c.qrUrl + "/" + nip + "/" + uploadedInvoice.InvoicingDate.Format("02-01-2006") + "/" + base64UrlHash,
		},
	)

	return nil
}
