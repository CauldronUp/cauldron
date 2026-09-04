# Pinecone

Emulates the Pinecone API (2026-04), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Pinecone is two APIs behind two different base URLs, and the response to "create an index" is mostly an address for somewhere else. The control plane (api.pinecone.io) creates and lists indexes; every actual operation an application performs, query, upsert, fetch, delete vectors, lives at the index's own host, which you cannot know until the control plane hands it to you. Pinecone's own data-plane spec is explicit about this: its server URL template defaults the `index_host` variable to the literal string `"unknown"`, so a client generated from that document and left unconfigured posts its first query to `https://unknown/query`.

Readiness is also two fields that can openly disagree -- Pinecone's own example shows `{"ready": true, "state": "ScalingUpPodSize"}`, ready while not in the state literally called Ready, out of nine possible states. Code waiting for `state == "Ready"` waits through scaling the index is already serving through, and code that only checks `ready` proceeds during `Terminating` too. A delete answers 202, not a clean removal, and `Terminating` is one of the nine states, so a deleted index keeps answering `describe` for a while afterward.

The error vocabulary is borrowed wholesale from gRPC -- codes like `UNKNOWN`, `DEADLINE_EXCEEDED`, `FAILED_PRECONDITION` sit alongside a few HTTP-flavored ones bolted on the end, and `OK` is a documented possible error code. The data plane itself is not served here, deliberately -- serving `/query` on the same host as the control plane would teach the exact mistake this Recipe exists to warn about, so what is modelled instead is the claim that produces it: an index's host is always somewhere else.

## Sources

- Documentation: https://docs.pinecone.io/reference/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pinecone     # run it
cauldron verify pinecone -v # check every claim
```
