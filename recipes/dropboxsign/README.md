# Dropbox Sign

Emulates the Dropbox Sign API (v3), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

There's no status field on a Dropbox Sign signature request -- just three independent booleans (is_complete, is_declined, has_error) that can combine in any way, including declined and errored at once. is_complete: false alone covers at least four different real situations: nobody has signed yet, somebody declined, something broke, or two of three signed and the third hasn't gotten to it -- so a dashboard that renders complete-or-pending shows a declined contract as pending forever.

The webhook contract is unusually literal: Dropbox Sign's own spec requires the callback endpoint to answer with the exact text "Hello API Event Received" as a plain-text body, not just a 200 -- every other provider in this collection reads your status code, this one reads your prose (Cauldron doesn't enforce this, since doing so would change how every Recipe's deliveries are judged, and it's recorded here rather than modelled). The signer who actually executes the document also need not be the person you invited: a signature carries reassigned_by, reassignment_reason and reassigned_from when someone hands their signing task to another person, so the name on the executed document can be someone the original request never named.

Test requests also appear in the same listing as real ones with no way to filter them out by the listing parameters offered, so a count of contracts sent includes every one a developer made while wiring the integration up.

## Sources

- Documentation: https://developers.hellosign.com/api/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve dropboxsign     # run it
cauldron verify dropboxsign -v # check every claim
```
