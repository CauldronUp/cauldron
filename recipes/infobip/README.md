# Infobip

Emulates the Infobip API (2), for local development and tests.

**8 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

It completes a five-way comparison the SMS Recipes have
been building. One request to many recipients gives back: nothing shared at all
from Plivo, one bare id from Bandwidth and Sinch, N separate requests from
Twilio -- and from Infobip alone, one entry per recipient in the first response,
each already carrying a structured status, so the answer a client gets first is
the one that can already disagree with itself.

## Sources

- Documentation: https://www.infobip.com/docs/sms/sms-over-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve infobip     # run it
cauldron verify infobip -v # check every claim
```
