package ksef

import (
	"fmt"

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
