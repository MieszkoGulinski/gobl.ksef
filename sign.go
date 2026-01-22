package ksef

import (
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/head"
)

// Sign attaches KSeF identification stamps (KSeF number, hash and QR URL) to the envelope.
func Sign(env *gobl.Envelope, qrURL, ksefNumber, invoiceHash string) error {
	if env == nil {
		return fmt.Errorf("envelope is nil")
	}

	env.Head.AddStamp(
		&head.Stamp{
			Provider: favat.StampKSEFNumber,
			Value:    ksefNumber,
		},
	)
	env.Head.AddStamp(
		&head.Stamp{
			Provider: favat.StampHash,
			Value:    invoiceHash,
		},
	)
	env.Head.AddStamp(
		&head.Stamp{
			Provider: favat.StampQR,
			Value:    qrURL,
		},
	)

	return nil
}
