# Documenso

Emulates the Documenso API (v1), for local development and tests.

**15 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A Documenso document and its recipients carry separate statuses answering different questions -- a document is PENDING while anyone hasn't signed yet, while each recipient is independently SIGNED, VIEWED, SENT or OPENED, so reading the document-level status to check whether one specific person signed gives you the answer for everybody at once. Sending a document also isn't signing it: a document moves to PENDING the moment it goes out and can sit there indefinitely waiting on humans.

A recipient's signing-URL token is not their id -- building a signing link from the id produces a URL that resolves to nothing. And recipients only get a sequential signing order when the document explicitly declares one; without it, everyone is emailed at once, so a workflow assuming sequential signing sends the second signer their link before the first has done anything.

Documenso is open source and genuinely self-hostable, unlike most providers here, but the cloud service has no sandbox -- sending a document for signature emails real people and produces a legal artefact that can't be unsent.

## Sources

- Documentation: https://docs.documenso.com/developers/public-api
- Machine-readable description: https://app.documenso.com/api/v1/openapi.json, last checked 2026-08-31
  `cauldron drift documenso` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve documenso     # run it
cauldron verify documenso -v # check every claim
```
