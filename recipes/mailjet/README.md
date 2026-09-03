# Mailjet

Emulates the Mailjet API (v3), for local development and tests.

**8 conformance cases, 3 checked against the live API on 2026-09-01.**

## Sources

- Documentation: https://dev.mailjet.com/email/reference/overview/introduction/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mailjet     # run it
cauldron verify mailjet -v # check every claim
```
