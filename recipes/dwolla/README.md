# dwolla

Emulates the dwolla API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A Dwolla create answers 201 with no body at all -- not an empty object, nothing -- so the identifier of what you just made lives only in the Location header, and a client that calls .json() on the response throws. Dwolla is also HAL throughout: resources carry no id field, only _links.self.href, so code reading customer.id gets undefined on every object and has to extract the id from the URL's last segment instead.

Micro-deposit verification has exactly three attempts before the funding source becomes permanently unverifiable -- not rate-limited, not retryable tomorrow, done -- and the only fix is removing it and adding it again, so a client that retries automatically on failure burns all three attempts in a second. Amounts are also decimal strings nested inside an object (amount.value), so reading transfer.amount gets an object and parsing it as a number gives NaN.

## Sources

- Documentation: https://developers.dwolla.com/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve dwolla     # run it
cauldron verify dwolla -v # check every claim
```
