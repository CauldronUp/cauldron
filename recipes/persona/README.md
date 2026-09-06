# Persona

Emulates the Persona API (2023-01-05), for local development and tests.

**14 conformance cases, 2 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs an account this Recipe cannot invent. The credential check itself was verified directly against api.withpersona.com on 2026-09-05.

## What this Recipe found

No Authorization header, a fictitious bearer token, an unrouted path, and a wrong method all answer the identical `{"errors":[{"title":"Must be authenticated to access this endpoint"}]}`, checked live: Persona checks the credential before anything else about the request, and does not distinguish an absent one from a wrong one. That contradicts what this file assumed before checking -- the message was written as "Invalid API key" -- and it also does not match Persona's own published error reference, which documents this exact failure as `{"title": "Unauthorized", "details": "An invalid API key was provided"}`. The reference's own general shape (plural `details`, no `code` field at all) was still wrong here too: the Recipe had been modelling a singular `detail` and a machine `code` neither one actually carries.

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
