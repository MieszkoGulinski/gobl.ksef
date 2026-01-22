//go:build xsdvalidate

package test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	xsdvalidate "github.com/terminalstatic/go-xsd-validate"
)

// ValidateAgainstFA3Schema validates the given data against the FA3 schema.
func ValidateAgainstFA3Schema(t *testing.T, data []byte) {
	err := xsdvalidate.Init()
	require.NoError(t, err)
	t.Cleanup(xsdvalidate.Cleanup)

	// Use file path instead of memory to allow relative imports to resolve
	schemaPath := filepath.Join(GetSchemaPath(), "FA3.xsd")
	xsdhandler, err := xsdvalidate.NewXsdHandlerUrl(schemaPath, xsdvalidate.ParsErrVerbose)
	require.NoError(t, err)
	t.Cleanup(xsdhandler.Free)

	validation := xsdhandler.ValidateMem(data, xsdvalidate.ParsErrDefault)
	assert.Nil(t, validation)
}






