# Persona

Emulates the Persona API (2023-01-05), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Persona speaks JSON:API, so nothing is where an ordinary REST client expects: the status lives at `data.attributes.status`, the id at `data.id`, and there is no top-level status or reference id at all -- a client that reads the object it was handed directly finds nothing. The verifications that actually explain a decision are not on the inquiry either; they arrive in a separate top-level `included` array that has to be matched back by type and id, so finding out why an inquiry was declined means walking a second collection.

`completed` on an inquiry means the person finished the flow, not that they passed -- whether they passed is a later value on the same status field, and code treating completion as a decision lets everyone through. `needs_review` is neither pass nor fail either; a human has to look, and it can sit there for hours, so treating it as failure turns away real customers and treating it as success is worse. `reference-id` (yours) and `id` (theirs) are also easy to mix up -- only one of them means anything on your own side, so storing the wrong one loses the link back to your own record entirely.

## Sources

- Documentation: https://docs.withpersona.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve persona     # run it
cauldron verify persona -v # check every claim
```
