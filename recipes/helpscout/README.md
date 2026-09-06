# Help Scout

Emulates the Help Scout API (v2), for local development and tests.

**21 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real mailbox, because a reply sends a real email. The refusal cases were struck live, unauthenticated, against api.helpscout.net.

## What this Recipe found

A thread of type `"note"` is internal and never shown to the customer, but it sits in the same array as the replies that were sent. Code that renders every thread into a customer-facing view publishes whatever a support agent privately wrote about that customer -- and the Recipe's own header calls this the most expensive mistake available in the API, one a hand-rolled fake would never surface because a fake only ever returns the threads somebody thought to write.

The response is HAL: collections live under `_embedded` keyed by resource name, paging under a `page` object, and the next-page cursor as a URL inside `_links`. A client reaching for `data`, `results`, or a bare top-level array finds nothing at any of them.

`status` has four values, and closed is only one of them: `"pending"` is a conversation someone has picked up but not answered, and `"spam"` still belongs to the customer and still comes back from a plain list. Treating anything not closed as still open miscounts both.

The live probe found authentication failures live on a completely different layer from every other error in this API: a missing or invalid token answers a flat OAuth2 object (`{"error":"...","error_description":"..."}`, RFC 6750), not the `_embedded.errors` HAL shape this file's other failures use and had wrongly declared for this one too. A missing token and an invalid one also carry different codes and sentences, and Help Scout checks the credential before it looks at the path or the method at all.

## Sources

- Documentation: https://developer.helpscout.com/mailbox-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve helpscout     # run it
cauldron verify helpscout -v # check every claim
```
