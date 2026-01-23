#!/bin/bash
#
# Convert all JSON invoices to XML and validate against FA3 schema
#
# Usage:
#   ./scripts/convert-validate.sh           # Convert and validate (requires libxml2)
#   ./scripts/convert-validate.sh --no-validate  # Convert only, skip validation
#

set -e

cd "$(dirname "$0")/.."

if [[ "$1" == "--no-validate" ]]; then
    echo "Converting all invoices (without XSD validation)..."
    go test ./test -run TestConvertAndValidateAll -v
else
    echo "Converting and validating all invoices against FA3 schema..."
    echo "(requires libxml2 to be installed)"
    echo ""
    go test -tags xsdvalidate ./test -run TestConvertAndValidateAll -v
fi


