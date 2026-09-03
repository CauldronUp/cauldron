# Api2Pdf

Emulates the Api2Pdf API (v2), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its retention is **24 hours** by its own FAQ, quoted
beside PDFMonkey's one hour so the comparison is between two vendors' words
rather than two measurements.

## Sources

- Documentation: https://www.api2pdf.com/documentation/v2
- Machine-readable description: https://v2.api2pdf.com/swagger/v2/swagger.json, last checked 2026-09-01
  `cauldron drift api2pdf` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve api2pdf     # run it
cauldron verify api2pdf -v # check every claim
```
