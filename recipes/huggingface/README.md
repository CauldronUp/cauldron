# Hugging Face

Emulates the Hugging Face API (hub), for local development and tests.

**21 conformance cases, 18 checked against the live API on 2026-08-31.** The unchecked one is the paging case: it sends the two parameter names this Recipe declares, which is a claim read from the provider's own description rather than struck against it.

## What this Recipe found

Face's, where **a model that does not exist is a 401.**
`{"error": "Invalid username or password."}` with a `WWW-Authenticate` header,
so a typo reads as a credentials problem and an anonymous client cannot tell a
wrong name from a private one. Its listings carry both `_id`, a Mongo ObjectId,
and `id`, the real slug -- and `_id` opens nothing: fetching by it gives the
same 401 as a name somebody invented. An invalid `sort` answers 400 whose body
says `"✖ Invalid sort parameter"` with U+2716 while the `X-Error-Message`
header says the same sentence with a plain ASCII asterisk.

## Sources

- Documentation: https://huggingface.co/docs/hub/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve huggingface     # run it
cauldron verify huggingface -v # check every claim
```
