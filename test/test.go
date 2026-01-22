// Package test provides tools for testing the library
package test

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/invopop/gobl"
	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/bill"
)

// UpdateOut is a flag that when set to true will regenerate all XML output files
var UpdateOut bool

func init() {
	flag.BoolVar(&UpdateOut, "update", false, "update the XML output files in test/data/out")
}

// BuildFAVATFrom creates a KSeF FA_VAT document from a GOBL file in the `test/data` folder.
func BuildFAVATFrom(name string) (*ksef.Invoice, error) {
	env, err := LoadTestEnvelope(name)
	if err != nil {
		return nil, err
	}

	return ksef.BuildFavat(env)
}

// NewDocumentFrom creates a KSeF Document from a GOBL file in the `test/data` folder.
//
// Deprecated: use BuildFAVATFrom.
func NewDocumentFrom(name string) (*ksef.Invoice, error) {
	return BuildFAVATFrom(name)
}

// LoadTestInvoice returns a GOBL Invoice from a file in the `test/data` folder
func LoadTestInvoice(name string) (*bill.Invoice, error) {
	env, err := LoadTestEnvelope(name)
	if err != nil {
		return nil, err
	}

	return env.Extract().(*bill.Invoice), nil
}

// LoadTestEnvelope returns a GOBL Envelope from a file in the `test/data` folder.
// It handles both envelope and direct invoice formats.
func LoadTestEnvelope(name string) (*gobl.Envelope, error) {
	return loadAndEnvelope(name)
}

// BuildFAVATFromInvoice returns a KSeF FA_VAT document from a GOBL invoice.
func BuildFAVATFromInvoice(inv *bill.Invoice) (*ksef.Invoice, error) {
	env, err := gobl.Envelop(inv)
	if err != nil {
		return nil, err
	}

	return ksef.BuildFavat(env)
}

// GenerateKSeFFrom returns a KSeF Document from a GOBL Invoice.
//
// Deprecated: use BuildFAVATFromInvoice.
func GenerateKSeFFrom(inv *bill.Invoice) (*ksef.Invoice, error) {
	return BuildFAVATFromInvoice(inv)
}

// LoadOutputFile returns byte data from a file in the `test/data/out` folder
func LoadOutputFile(name string) ([]byte, error) {
	src, _ := os.Open(filepath.Join(GetOutPath(), name))

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// LoadSchemaFile returns byte data from a file in the `test/data/schema` folder
func LoadSchemaFile(name string) ([]byte, error) {
	src, _ := os.Open(filepath.Join(GetSchemaPath(), name))

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// loadAndEnvelope loads a JSON file and returns it as an envelope.
// It handles both envelope and direct invoice formats.
func loadAndEnvelope(name string) (*gobl.Envelope, error) {
	src, err := os.Open(filepath.Join(GetDataPath(), name))
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	data := buf.Bytes()

	// Try to parse as envelope first
	env := new(gobl.Envelope)
	if err := json.Unmarshal(data, env); err == nil && env.Document != nil {
		if err := env.Calculate(); err != nil {
			return nil, err
		}
		return env, nil
	}

	// Try to parse as invoice
	inv := new(bill.Invoice)
	if err := json.Unmarshal(data, inv); err != nil {
		return nil, err
	}

	// Wrap in envelope
	env, err = gobl.Envelop(inv)
	if err != nil {
		return nil, err
	}

	return env, nil
}

// GetSchemaPath returns the path to the `test/data/schema` folder
func GetSchemaPath() string {
	return filepath.Join(GetDataPath(), "schema")
}

// GetOutPath returns the path to the `test/data/out` folder
func GetOutPath() string {
	return filepath.Join(GetDataPath(), "out")
}

// GetDataPath returns the path to the `test/data` folder
func GetDataPath() string {
	return filepath.Join(GetTestPath(), "data")
}

// GetTestPath returns the path to the `test` folder
func GetTestPath() string {
	return filepath.Join(getRootFolder(), "test")
}

func getRootFolder() string {
	cwd, _ := os.Getwd()

	for !isRootFolder(cwd) {
		cwd = removeLastEntry(cwd)
	}

	return cwd
}

func isRootFolder(dir string) bool {
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if file.Name() == "go.mod" {
			return true
		}
	}

	return false
}

func removeLastEntry(dir string) string {
	lastEntry := "/" + filepath.Base(dir)
	i := strings.LastIndex(dir, lastEntry)
	return dir[:i]
}

// ValidateAgainstFA3Schema is implemented in build-tagged files:
// - with `-tags xsdvalidate`: runs schema validation via libxml2
// - without it: skips schema validation
