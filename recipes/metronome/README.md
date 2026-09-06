# Metronome

Emulates the Metronome API (v1), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The endpoint that decides what a customer gets charged answers 200 and nothing else. `POST /v1/ingest`'s entire published response is `"200: Success"` -- no schema, no body, no per-event result -- and Metronome's own docs say duplicates are silently detected and dropped inside a 34-day window. An event that was billed and an event that was silently discarded produce the identical answer, and there is nothing in the response to compare, log, or alert on. The only endpoint that could tell you is explicitly rate-limited for sampling only, with a warning not to use it to check every event -- so whether any specific invoice is right is not a question this API can answer at all.

The 34-day dedup window cuts both ways: it exists to make backfills safe, but it also means a corrected event resent under the same transaction id is treated as a duplicate and silently dropped for a month, so a fix for a wrong number gets discarded and the wrong number stands. Backdating an event up to 34 days back immediately changes what a customer owes right now, with no confirmation step. An event can also be accepted as a 200 and belong to nobody -- `matched_customer` and `matched_billable_metrics` are simply absent when nothing matched, so the event is stored and searchable but bills nothing, and there is no way from the ingest response alone to know that happened.

The whole error schema, everywhere, is `{"message": string}` -- no code, no type, no field -- so every retry decision this API supports is a decision about a sentence rather than a structured reason.

## Sources

- Documentation: https://docs.metronome.com/api
- Machine-readable description: https://docs.metronome.com/openapi.json, last checked 2026-09-05
  `cauldron drift metronome` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve metronome     # run it
cauldron verify metronome -v # check every claim
```
