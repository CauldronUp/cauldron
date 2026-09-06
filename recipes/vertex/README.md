# Vertex

Emulates the Vertex API (v2), for local development and tests.

**4 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The header for this Recipe is mostly a record of a dead end: `api.vertexcloud.com`, the host this Recipe was asked to check, does not resolve at all on any resolver tried, and the one server Vertex's own OpenAPI description names under that domain, `developer.vertexcloud.com`, is explicitly labelled a development environment and points to a dangling CNAME. There was no live host to probe, so every case here carries no verified date -- not because nobody tried, but because trying is recorded in prose instead of a request that could be sent.

What the description does establish, read rather than called: one endpoint, `/supplies`, decides between an estimate and a stored, reportable transaction purely by the value of one field, `saleMessageType` -- QUOTATION versus INVOICE versus a tax-only adjustment versus a reconciliation sync. And undoing a stored transaction is two separate mechanisms, not one: a delete that only flags a transaction as excluded from reporting rather than removing it, and a reversal that creates a brand new negating transaction rather than touching the original.

## Sources

- Documentation: https://developer.vertexinc.com/oseries
- Machine-readable description: https://dash.readme.com/api/v1/api-registry/1bto7bwmt0de2rl, last checked 2026-09-01
  `cauldron drift vertex` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve vertex     # run it
cauldron verify vertex -v # check every claim
```
