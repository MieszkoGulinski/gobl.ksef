//go:build !xsdvalidate

package test

import "testing"

// ValidateAgainstFA3Schema skips schema validation unless built with `-tags xsdvalidate`.
func ValidateAgainstFA3Schema(t *testing.T, _ []byte) {
	t.Skip("FA3 XSD validation requires libxml2; run with `go test -tags xsdvalidate ./...`")
}








