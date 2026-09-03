# EPSS

Emulates the EPSS API (epss), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**One of the keys has a hyphen in it.** The envelope is
`{"status": "OK", "status-code": 200, ...}`, and `status-code` cannot be reached
with a dot in any language whose dot operator takes an identifier. In
JavaScript, `response["status-code"]` is 200 and `response.status-code` is
`response.status` minus a variable called `code` -- `NaN`, at runtime, with no
error and no warning. The field most likely to be read by code checking whether
the request worked is the one most likely to be read wrong.

**And the status is stated three times:** the HTTP status line, `status`, and
`status-code`, three sources for one fact in one response. The probability is a
string to nine decimal places -- `"0.999990000"`, quoted and padded, wrong for
every numeric comparison that does not parse it first. A CVE that is not a CVE
is a 200 with an empty array, exactly as an unknown but well-formed one is,
while omitting the parameter entirely is a 404: the malformed request succeeds
and the incomplete one fails. And failing makes the data private -- a 200 says
`"access": "public"` and a 404 says `"access": "private"`, so a field describing
what may be done with the data changes because the request did not work.

## Sources

- Documentation: https://www.first.org/epss/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve epss     # run it
cauldron verify epss -v # check every claim
```
