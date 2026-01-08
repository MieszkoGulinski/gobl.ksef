package api

import (
	"encoding/xml"
	"fmt"
	"net/url"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/head"
)

// UPO defines the XML structure for KSeF UPO
type UPO struct {
	KSeFNumber string `xml:"Dokument>NumerKSeFDokumentu"`
	KSeFHash   string `xml:"Dokument>SkrotDokumentu"`
}

// Sign reads the UPO file and adds the QR code values to the envelope
func Sign(env *gobl.Envelope, upoBytes []byte, c *Client) error {
	upo := new(UPO)
	if err := xml.Unmarshal(upoBytes, upo); err != nil {
		return fmt.Errorf("parsing input as UPO: %w", err)
	}

	env.Head.AddStamp(
		&head.Stamp{
			Provider: favat.StampID,
			Value:    upo.KSeFNumber,
		},
	)
	env.Head.AddStamp(
		&head.Stamp{
			Provider: favat.StampHash,
			Value:    upo.KSeFHash,
		},
	)
	env.Head.AddStamp(
		&head.Stamp{
			Provider: favat.StampQR,
			Value:    c.url + "/web/verify/" + upo.KSeFNumber + "/" + url.QueryEscape(upo.KSeFHash), // TODO check if this is correct
		},
	)

	return nil
}
