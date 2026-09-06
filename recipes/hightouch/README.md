# hightouch

Emulates the hightouch API (v1), for local development and tests.

**19 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real workspace; the refusal cases were struck live, unauthenticated, against api.hightouch.com.

## What this Recipe found

Reverse ETL has a failure mode no CRUD API has: the job can succeed while the data never arrives. A sync run reports rows queried, planned, attempted and rejected as separate counts, and its own `status` is only about whether the machinery ran, not whether the rows landed -- a run that attempts four thousand rows and has every single one rejected by the destination still finishes with `status: "success"`. The dashboard is green. Nobody's CRM has the data.

Row-level rejection reasons live behind a completely different endpoint from the run summary, and almost nobody fetches them. `"warning"` is also a real status, distinct from both success and failure -- it means some rows landed and some didn't, and code that treats anything short of `"failed"` as fine treats a partial delivery as a complete one. A sync can also be stopped two different ways that look identical from outside: disabled by a person, or paused by Hightouch after repeated failures -- and the fix for each is different.

The live probe found this file had the message and the detail backwards: the real API sends a constant top-level `"message":"Authentication error"` with the specific reason in a separate `details` field, not the single varying sentence this file declared. A missing header and a wrong token share the message and differ only in that detail, and Hightouch checks the credential before it looks at the path or the method at all.

## Sources

- Documentation: https://hightouch.com/docs/api-reference/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve hightouch     # run it
cauldron verify hightouch -v # check every claim
```
