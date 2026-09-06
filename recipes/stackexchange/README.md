# Stack Exchange

Emulates the Stack Exchange API (v2.3), for local development and tests.

**9 conformance cases, 8 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Exchange's, where **every failure is 400 and one of them carries
a 404 in the body.** A path that does not exist answers `400` on the wire with
`{"error_id": 404, "error_name": "no_method"}` in the payload, so a client
switching on the status cannot tell "you asked wrongly" from "there is nothing
there", and one switching on `error_id` disagrees with its own transport.

**And `error_message` is sometimes a sentence and sometimes just a parameter
name.** A bad site says ``No site found for name `nosuchsite` `` -- Markdown
backticks in a response that is not Markdown -- while a bad page size says,
in full, `"pagesize"`, and a bad id says `"ids"`. Anything showing
`error_message` to a person shows them the word "ids". The rate-limit budget
rides in the body of every answer, `quota_max` and `quota_remaining` beside the
results, so the payload differs on every request even when the data does not.
There is no total anywhere -- only `has_more`, so a pager cannot be built from
it. Every timestamp is a Unix epoch integer, and `last_activity_date` and
`last_edit_date` are the same number on this question. And every item carries
its own `content_license`.

## Sources

- Documentation: https://api.stackexchange.com/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve stackexchange     # run it
cauldron verify stackexchange -v # check every claim
```
