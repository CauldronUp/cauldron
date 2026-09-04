# Help Scout

Emulates the Help Scout API (v2), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A thread of type `"note"` is internal and never shown to the customer, but it sits in the same array as the replies that were sent. Code that renders every thread into a customer-facing view publishes whatever a support agent privately wrote about that customer -- and the Recipe's own header calls this the most expensive mistake available in the API, one a hand-rolled fake would never surface because a fake only ever returns the threads somebody thought to write.

The response is HAL: collections live under `_embedded` keyed by resource name, paging under a `page` object, and the next-page cursor as a URL inside `_links`. A client reaching for `data`, `results`, or a bare top-level array finds nothing at any of them.

`status` has four values, and closed is only one of them: `"pending"` is a conversation someone has picked up but not answered, and `"spam"` still belongs to the customer and still comes back from a plain list. Treating anything not closed as still open miscounts both.

## Sources

- Documentation: https://developer.helpscout.com/mailbox-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve helpscout     # run it
cauldron verify helpscout -v # check every claim
```
