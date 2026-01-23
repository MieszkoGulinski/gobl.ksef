// Package ksef implements the conversion from GOBL to FA_VAT XML
package ksef

import (
	"time"
)

// KSEF schema constants
const (
	systemCode    = "FA (3)"
	formCode      = "FA"
	schemaVersion = "1-0E"
	formVariant   = 3
	systemInfo    = "Invopop"
)

// Header defines the XML structure for KSeF header
type Header struct {
	FormCode     *FormCode `xml:"KodFormularza"`
	FormVariant  int       `xml:"WariantFormularza"`
	CreationDate string    `xml:"DataWytworzeniaFa"`
	SystemInfo   string    `xml:"SystemInfo,omitempty"`
}

// FormCode defines the XML structure for KSeF schema versioning
type FormCode struct {
	SystemCode    string `xml:"kodSystemowy,attr"`
	SchemaVersion string `xml:"wersjaSchemy,attr"`
	FormCode      string `xml:",chardata"`
}

// NewFavatHeader gets header data from GOBL invoice
func NewFavatHeader() *Header {
	header := &Header{
		FormCode: &FormCode{
			SystemCode:    systemCode,
			SchemaVersion: schemaVersion,
			FormCode:      formCode,
		},
		FormVariant:  formVariant,
		CreationDate: formatGenerationDate(time.Now()),
		SystemInfo:   systemInfo,
	}

	return header
}

func formatGenerationDate(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05Z")
}
