package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/stretchr/testify/require"
)

// TestKSeF tests all GOBL JSON files in test/data.
//
// Without --update flag:
//   - Reads existing XML files from test/data/out
//   - Validates them against the FA3 schema
//
// With --update flag:
//   - Converts GOBL JSON to KSeF XML
//   - Saves XML to test/data/out
//   - Validates against the FA3 schema
//
// Run without XSD validation:
//
//	go test ./test -v
//	go test ./test --update -v
//
// Run with XSD validation (requires libxml2):
//
//	go test -tags xsdvalidate ./test -v
//	go test -tags xsdvalidate ./test --update -v
func TestKSeF(t *testing.T) {
	dataPath := GetDataPath()
	outPath := GetOutPath()

	// Ensure output directory exists
	err := os.MkdirAll(outPath, 0755)
	require.NoError(t, err)

	// Read all JSON files in data directory
	entries, err := os.ReadDir(dataPath)
	require.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		name := entry.Name()
		t.Run(name, func(t *testing.T) {
			xmlName := strings.TrimSuffix(name, ".json") + ".xml"
			xmlPath := filepath.Join(outPath, xmlName)

			var xmlData []byte

			if UpdateOut {
				// Update mode: convert GOBL JSON to KSeF XML
				t.Logf("Converting %s to XML...", name)

				env, err := LoadTestEnvelope(name)
				require.NoError(t, err, "failed to load envelope")

				doc, err := ksef.BuildFavat(env)
				require.NoError(t, err, "failed to build FA_VAT document")

				xmlData, err = doc.Bytes()
				require.NoError(t, err, "failed to generate XML bytes")

				// Save to output directory
				err = os.WriteFile(xmlPath, xmlData, 0644)
				require.NoError(t, err, "failed to write XML file")

				t.Logf("Generated %d bytes of XML → %s", len(xmlData), xmlName)
			} else {
				// Validate mode: read existing XML
				xmlData, err = os.ReadFile(xmlPath)
				require.NoError(t, err, "failed to read XML file (run with --update to generate)")

				t.Logf("Validating existing XML: %s (%d bytes)", xmlName, len(xmlData))
			}

			// Validate against FA3 schema
			ValidateAgainstFA3Schema(t, xmlData)
		})
	}
}
