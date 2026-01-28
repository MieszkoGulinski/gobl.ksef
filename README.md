# GOBL ↔ KSeF Conversion

Convert GOBL to the Polish FA_VAT XML format (and, in future, back from KSeF into GOBL).

The main conversion entrypoints are `ksef.BuildFAVAT` (model) and `ksef.MarshalFAVAT` (XML bytes).

Copyright [Invopop Ltd.](https://invopop.com) 2023. Released publicly under the [Apache License Version 2.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). In order to accept contributions to this library we will require transferring copyrights to Invopop Ltd.

## Project Development Objectives

The following list the steps to follow through on in order to accomplish the goal of using GOBL to submit electronic invoices to the Polish authorities:

1. Add the PL (`pl`) tax regime to [GOBL](https://github.com/invopop/gobl). Figure out local taxes, tax ID validation rules, and any "extensions" that may be required to be defined in GOBL, and send in a PR. For examples of existing regimes, see the [regimes](https://github.com/invopop/gobl/tree/main/regimes) directory. Key Concerns:
   - Basic B2B invoices support.
   - Tax ID validation as per local rules.
   - Support for "simplified" invoices.
   - Requirements for credit-notes or "rectified" invoices and the correction options definition for the tax regime.
   - Any additional fields that need to be validated, like payment terms.
2. Convert GOBL into FA_VAT format in library. A couple of good examples: [gobl.cfdi for Mexico](https://github.com/invopop/gobl.cfdi) and [gobl.verifactu for Spain](https://github.com/invopop/gobl.verifactu). Library would just be able to run tests in the first version.
3. Build a CLI (copy from gobl.cfdi and gobl.verifactu projects) to convert GOBL JSON documents into FA_VAT XML.
4. Build a second part of this project that allows documents to be sent directly to the KSeF. A partial example of this can be found in the [gobl.ticketbai project](https://github.com/invopop/gobl.ticketbai/tree/refactor/internal/gateways). It'd probably be useful to be able to upload via the CLI too.

## Testing

The test suite automatically discovers and tests all GOBL JSON files in `test/data/`.

### Running Tests

**Validate existing XML files against the schema:**
```bash
go test ./test -v
```

**Convert GOBL to KSeF XML and validate:**
```bash
go test ./test --update -v
```

**With XSD schema validation (requires libxml2):**
```bash
# Using the helper script (sets LD_LIBRARY_PATH automatically)
./test/test.sh -v
./test/test.sh --update -v

# Or manually
LD_LIBRARY_PATH=/home/linuxbrew/.linuxbrew/opt/libxml2/lib:$LD_LIBRARY_PATH go test -tags xsdvalidate ./test -v
```

### Test Data

- **Input**: GOBL JSON files in `test/data/*.json`
- **Output**: KSeF XML files in `test/data/out/*.xml`
- **Schema**: FA3 XSD and dependencies in `test/data/schema/`

## Unsupported fields

See [unsupported-fields.md](unsupported-fields.md) for the list of unsupported fields.

## FA_VAT documentation

FA_VAT is the Polish electronic invoice format. The format uses XML.

- [XML schema](https://github.com/CIRFMF/ksef-docs/blob/main/faktury/schemy/FA/schemat_FA(3)_v1-0E.xsd) for V3 (description of fields is in Polish)
- [Types definition](https://raw.githubusercontent.com/CIRFMF/ksef-docs/refs/heads/main/faktury/schemy/FA/bazowe/ElementarneTypyDanych_v10-0E.xsd) (description of fields is in Polish) - we have to open it as raw, as [the original link](https://github.com/CIRFMF/ksef-docs/blob/main/faktury/schemy/FA/bazowe/StrukturyDanych_v10-0E.xsd) does not add newlines
- [Complex types definition](https://raw.githubusercontent.com/CIRFMF/ksef-docs/refs/heads/main/faktury/schemy/FA/bazowe/StrukturyDanych_v10-0E.xsd) (description of fields is in Polish) - we have to open it as raw, as [the original link](https://github.com/CIRFMF/ksef-docs/blob/main/faktury/schemy/FA/bazowe/StrukturyDanych_v10-0E.xsd) does not add newlines

## Parsing
The parsing is in MVP and some cases are not entirely handle. They will be handle little by little once we start facing them. There might be cases where the invoices from the tax agency don't come with lines, which makes it difficult to parse on our sie as we would need to fake lines.

## KSeF API

KSeF is the Polish system for submitting electronic invoices to the Polish authorities.

Useful links:

- [National e-Invoice System](https://ksef.mf.gov.pl/) - for details on system in general (English translation available - language picker is in the top right corner)
- [KSeF Test Zone](https://ksef-test.mf.gov.pl/) - as above, but for testing
- [API documentation](https://ksef-test.mf.gov.pl/docs/v2/index.html) for the test environment (in Polish)

KSeF provide three environments:

1.  [Test Environment](https://ksef-test.mf.gov.pl/) for application development with fictitious data.
2.  [Pre-production "demo"](https://ksef-demo.mf.gov.pl/) area with production data, but not officially declared.
3.  [Production](https://ksef.mf.gov.pl)

A translation of the Interface Specification 1.5 is available in the [docs](./docs) folder.

OpenAPI documentation is available for three specific interfaces:

1. Batches ([test openapi 'batch' spec](https://ksef-test.mf.gov.pl/openapi/gtw/svc/api/KSeF-batch.yaml)) - for sending multiple documents at the same time.
2. Common ([test openapi 'common' spec](https://ksef-test.mf.gov.pl/openapi/gtw/svc/api/KSeF-common.yaml)) - general operations that don't require authentication.
3. Interactive ([test openapi 'online' spec](https://ksef-test.mf.gov.pl/openapi/gtw/svc/api/KSeF-online.yaml)) - sending a single document in each request.

## Authentication

See [authentication.md](./authentication.md).
